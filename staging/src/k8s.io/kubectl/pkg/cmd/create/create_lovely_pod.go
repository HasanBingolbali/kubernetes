package create

import (
	"context"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/cli-runtime/pkg/genericiooptions"
	coreclient "k8s.io/client-go/kubernetes/typed/core/v1"
	cmdutil "k8s.io/kubectl/pkg/cmd/util"
	"k8s.io/kubectl/pkg/scheme"
	"k8s.io/kubectl/pkg/util/i18n"
)

type RunOptions struct {
	cmdutil.OverrideOptions
	PrintFlags *genericclioptions.PrintFlags
	PrintObj   func(obj runtime.Object) error
	Client     *coreclient.CoreV1Client
	genericiooptions.IOStreams
}

func NewRunOptions(ioStreams genericiooptions.IOStreams) *RunOptions {
	return &RunOptions{
		PrintFlags: genericclioptions.NewPrintFlags("created").WithTypeSetter(scheme.Scheme),
		IOStreams:  ioStreams,
	}
}

func NewCmdRun(f cmdutil.Factory, streams genericiooptions.IOStreams) *cobra.Command {
	o := NewRunOptions(streams)

	cmd := &cobra.Command{
		Use:                   "lovely-pod",
		DisableFlagsInUseLine: true,
		Short:                 i18n.T("Run a particular image on the cluster"),
		Long:                  "runLong",
		Example:               "runExample",
		Run: func(cmd *cobra.Command, args []string) {
			cmdutil.CheckErr(o.Run(f, cmd, args))
		},
	}

	return cmd
}

func (o *RunOptions) Run(f cmdutil.Factory, cmd *cobra.Command, args []string) error {
	var err error
	restConfig, err := f.ToRESTConfig()
	if err != nil {
		return err
	}
	o.Client, err = coreclient.NewForConfig(restConfig)
	if err != nil {
		return err
	}
	printer, err := o.PrintFlags.ToPrinter()
	if err != nil {
		return err
	}
	o.PrintObj = func(obj runtime.Object) error {
		return printer.PrintObj(obj, o.Out)
	}

	pod, err := o.Client.Pods("default").Create(context.TODO(), &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "iloveyou",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx",
				},
			},
		},
	}, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return o.PrintObj(pod)
}
