package cmd

import (
	"context"
	"strconv"

	cliutil "github.com/k3d-io/k3d/v5/cmd/util"

	"github.com/docker/go-connections/nat"
	k3dCluster "github.com/k3d-io/k3d/v5/pkg/client"
	"github.com/k3d-io/k3d/v5/pkg/config"
	conf "github.com/k3d-io/k3d/v5/pkg/config/v1alpha4"
	l "github.com/k3d-io/k3d/v5/pkg/logger"
	"github.com/k3d-io/k3d/v5/pkg/runtimes"
	k3d "github.com/k3d-io/k3d/v5/pkg/types"
	"github.com/spf13/viper"
)

var configFile string
var filters []string

var (
	cfgViper = viper.New()
	ppViper  = viper.New()
)

func NewCmdClusterCreate() (string, error) {
	var (
		err       error
		exposeAPI *k3d.ExposureOpts
	)
	simpleCfg, err := config.SimpleConfigFromViper(cfgViper)
	exposeAPI = &k3d.ExposureOpts{
		PortMapping: nat.PortMapping{
			Binding: nat.PortBinding{
				HostIP:   simpleCfg.ExposeAPI.HostIP,
				HostPort: simpleCfg.ExposeAPI.HostPort,
			},
		},
		Host: simpleCfg.ExposeAPI.Host,
	}
	if len(exposeAPI.Binding.HostPort) == 0 {
		var freePort string
		port, err := cliutil.GetFreePort()
		freePort = strconv.Itoa(port)
		if err != nil || port == 0 {
			l.Log().Warnf("Failed to get random free port: %+v", err)
			l.Log().Warnf("Falling back to internal port %s (may be blocked though)...", k3d.DefaultAPIPort)
			freePort = k3d.DefaultAPIPort
		}
		exposeAPI.Binding.HostPort = freePort
	}
	simpleCfg.ExposeAPI = conf.SimpleExposureOpts{
		Host:     exposeAPI.Host,
		HostIP:   exposeAPI.Binding.HostIP,
		HostPort: exposeAPI.Binding.HostPort,
	}

	simpleCfg.Name = "tanuu"
	simpleCfg.Servers = 1
	simpleCfg.Agents = 0
	simpleCfg.Image = "docker.io/rancher/k3s:v1.23.8-k3s1"
	if demo {
		simpleCfg.Volumes = append(simpleCfg.Volumes, conf.VolumeWithNodeFilters{
			Volume:      viper.GetString("localrepo") + "/crds:/var/lib/rancher/k3s/server/manifests/crds",
			NodeFilters: filters,
		})
		simpleCfg.Volumes = append(simpleCfg.Volumes, conf.VolumeWithNodeFilters{
			Volume:      viper.GetString("localrepo") + "/demo:/var/lib/rancher/k3s/server/manifests/demo",
			NodeFilters: filters,
		})
		simpleCfg.Volumes = append(simpleCfg.Volumes, conf.VolumeWithNodeFilters{
			Volume:      viper.GetString("localrepo") + "/bootstrap:/var/lib/rancher/k3s/server/manifests/initial",
			NodeFilters: filters,
		})
	}
	if viper.GetBool("aws") {
		simpleCfg.Volumes = append(simpleCfg.Volumes, conf.VolumeWithNodeFilters{
			Volume:      viper.GetString("localrepo") + "/mgmtCluster/aws:/var/lib/rancher/k3s/server/manifests/aws",
			NodeFilters: filters,
		})
		simpleCfg.Volumes = append(simpleCfg.Volumes, conf.VolumeWithNodeFilters{
			Volume:      viper.GetString("localrepo") + "/mgmtCluster/base:/var/lib/rancher/k3s/server/manifests/base",
			NodeFilters: filters,
		})
	}
	if viper.GetBool("azure") {
		simpleCfg.Volumes = append(simpleCfg.Volumes, conf.VolumeWithNodeFilters{
			Volume:      viper.GetString("localrepo") + "/mgmtCluster/azure:/var/lib/rancher/k3s/server/manifests/azure",
			NodeFilters: filters,
		})
		simpleCfg.Volumes = append(simpleCfg.Volumes, conf.VolumeWithNodeFilters{
			Volume:      viper.GetString("localrepo") + "/mgmtCluster/base:/var/lib/rancher/k3s/server/manifests/base",
			NodeFilters: filters,
		})
	}
	filters = append(filters, "server:0")
	simpleCfg.Ports = append(simpleCfg.Ports, conf.PortWithNodeFilters{
		Port:        viper.GetString("port1") + ":30081",
		NodeFilters: filters,
	})
	simpleCfg.Ports = append(simpleCfg.Ports, conf.PortWithNodeFilters{
		Port:        viper.GetString("port2") + ":30082",
		NodeFilters: filters,
	})
	simpleCfg.Ports = append(simpleCfg.Ports, conf.PortWithNodeFilters{
		Port:        viper.GetString("port3") + ":30083",
		NodeFilters: filters,
	})
	simpleCfg.Ports = append(simpleCfg.Ports, conf.PortWithNodeFilters{
		Port:        viper.GetString("port4") + ":30084",
		NodeFilters: filters,
	})
	simpleCfg.Ports = append(simpleCfg.Ports, conf.PortWithNodeFilters{
		Port:        viper.GetString("port5") + ":30085",
		NodeFilters: filters,
	})

	clusterConfig, err := config.TransformSimpleToClusterConfig(context.TODO(), runtimes.SelectedRuntime, simpleCfg)
	if err != nil {
		l.Log().Fatalln(err)
	}
	clusterConfig, err = config.ProcessClusterConfig(*clusterConfig)
	if err != nil {
		l.Log().Fatalf("error processing cluster configuration: %v", err)
		return "", err
	}

	if err := config.ValidateClusterConfig(context.TODO(), runtimes.SelectedRuntime, *clusterConfig); err != nil {
		l.Log().Fatalln("Failed Cluster Configuration Validation: ", err)
	}

	if _, err := k3dCluster.ClusterGet(context.TODO(), runtimes.SelectedRuntime, &clusterConfig.Cluster); err == nil {
		l.Log().Fatalf("Failed to create cluster '%s' because a cluster with that name already exists", clusterConfig.Cluster.Name)
	}

	// create cluster
	clusterConfig.KubeconfigOpts.SwitchCurrentContext = true
	clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig = true
	if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig {
		l.Log().Debugln("'--kubeconfig-update-default set: enabling wait-for-server")
		clusterConfig.ClusterCreateOpts.WaitForServer = true
	}
	if err := k3dCluster.ClusterRun(context.TODO(), runtimes.SelectedRuntime, clusterConfig); err != nil {
		l.Log().Errorln(err)
		if simpleCfg.Options.K3dOptions.NoRollback {
			l.Log().Fatalln("Cluster creation FAILED, rollback deactivated.")
		}
		l.Log().Errorln("Failed to create cluster >>> Rolling Back")
		if err := k3dCluster.ClusterDelete(context.TODO(), runtimes.SelectedRuntime, &clusterConfig.Cluster, k3d.ClusterDeleteOpts{SkipRegistryCheck: true}); err != nil {
			l.Log().Errorln(err)
			l.Log().Fatalln("Cluster creation FAILED, also FAILED to rollback changes!")
		}
		l.Log().Fatalln("Cluster creation FAILED, all changes have been rolled back!")
	}
	l.Log().Infof("Cluster '%s' created successfully!", clusterConfig.Cluster.Name)

	kubecfg, err := k3dCluster.KubeconfigGetWrite(context.TODO(), runtimes.SelectedRuntime, &clusterConfig.Cluster, "", &k3dCluster.WriteKubeConfigOptions{})
	if err != nil {
		l.Log().Error(err)
	}
	return kubecfg, nil
}

func DeleteCluster() {
	simpleCfg, err := config.SimpleConfigFromViper(cfgViper)
	simpleCfg.Name = "tanuu"
	clusterConfig, err := config.TransformSimpleToClusterConfig(context.TODO(), runtimes.SelectedRuntime, simpleCfg)
	if err != nil {
		l.Log().Fatalln(err)
	}
	k3dCluster.ClusterDelete(context.TODO(), runtimes.SelectedRuntime, &clusterConfig.Cluster, k3d.ClusterDeleteOpts{})
}
