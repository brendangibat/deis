package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/brendangibat/deis/tests/dockercli"
	"github.com/brendangibat/deis/tests/etcdutils"
	"github.com/brendangibat/deis/tests/utils"

	"io/ioutil"
	"os"
)

func TestBuilder(t *testing.T) {
	var err error
	var errfile error
	setkeys := []string{
		"/deis/registry/protocol",
		"/deis/registry/host",
		"/deis/registry/port",
		"/deis/cache/host",
		"/deis/cache/port",
		"/deis/controller/protocol",
		"/deis/controller/host",
		"/deis/controller/port",
		"/deis/controller/builderKey",
	}
	setdir := []string{
		"/deis/controller",
		"/deis/cache",
		"/deis/database",
		"/deis/registry",
		"/deis/domains",
		"/deis/services",
	}
	setproxy := []byte("HTTP_PROXY=\nhttp_proxy=\n")

	tmpfile, errfile := ioutil.TempFile("/tmp", "deis-test-")
	if errfile != nil {
		t.Fatal(errfile)
	}
	ioutil.WriteFile(tmpfile.Name(), setproxy, 0644)
	defer os.Remove(tmpfile.Name())

	tag, etcdPort := utils.BuildTag(), utils.RandomPort()
	imageName := utils.ImagePrefix() + "builder" + ":" + tag
	etcdName := "deis-etcd-" + tag
	cli, stdout, stdoutPipe := dockercli.NewClient()
	dockercli.RunTestEtcd(t, etcdName, etcdPort)
	defer cli.CmdRm("-f", etcdName)
	handler := etcdutils.InitEtcd(setdir, setkeys, etcdPort)
	etcdutils.PublishEtcd(t, handler)
	host, port := utils.HostAddress(), utils.RandomPort()
	fmt.Printf("--- Run %s at %s:%s\n", imageName, host, port)
	name := "deis-builder-" + tag
	defer cli.CmdRm("-f", "-v", name)
	go func() {
		_ = cli.CmdRm("-f", "-v", name)
		err = dockercli.RunContainer(cli,
			"--name", name,
			"--rm",
			"-p", port+":22",
			"-e", "PORT=22",
			"-e", "HOST="+host,
			"-e", "ETCD_PORT="+etcdPort,
			"-e", "EXTERNAL_PORT="+port,
			"--privileged",
			"-v", tmpfile.Name()+":/etc/environment_proxy",
			imageName)
	}()
	dockercli.PrintToStdout(t, stdout, stdoutPipe, "deis-builder running")
	if err != nil {
		t.Fatal(err)
	}
	// FIXME: builder needs a few seconds to wake up here!
	// FIXME: Wait until etcd keys are published
	time.Sleep(5000 * time.Millisecond)
	dockercli.DeisServiceTest(t, name, port, "tcp")
	etcdutils.VerifyEtcdValue(t, "/deis/builder/host", host, etcdPort)
	etcdutils.VerifyEtcdValue(t, "/deis/builder/port", port, etcdPort)
}
