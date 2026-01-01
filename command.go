package main

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jmespath/go-jmespath"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

type AddOptions struct {
	Query     string
	EndPoint  string
	RawOutput bool
}

type DeleteOptions struct {
	Query     string
	EndPoint  string
	RawOutput bool
	SilenceId string
}

type ListOptions struct {
	Query     string
	EndPoint  string
	RawOutput bool
}

type ListWorkspacesOptions struct {
	Query     string
	RawOutput bool
}

func rootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "amp-silence",
		Short:         "",
		Args:          cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
		SilenceUsage:  true,
		SilenceErrors: false,
	}

	cmd.AddCommand(createVersionCmd())
	cmd.AddCommand(createAddCmd())
	cmd.AddCommand(createDeleteCommand())
	cmd.AddCommand(createListCmd())
	cmd.AddCommand(createListWorkspcesCmd())

	return cmd
}

func applyJMESPath(query string, data []byte, rawOutput bool) ([]byte, error) {
	var jsonData interface{}
	err := json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}
	result, err := jmespath.Search(query, jsonData)
	if err != nil {
		return nil, err
	}

	if rawOutput {
		switch t := result.(type) {
		case string:
			return []byte(t), nil
		case *string:
			return []byte(aws.ToString(t)), nil
		}
	}

	return json.Marshal(result)
}

func createVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "version",
		Short:         "",
		Args:          cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			if info, ok := debug.ReadBuildInfo(); ok {
				fmt.Printf("Version: %s\n", info.Main.Version)
				return nil
			} else {
				return errors.Errorf("could not read build info")
			}
		},
	}
	return cmd
}

func printResult(result []byte, query string, rawOutput bool) {
	if len(query) > 0 {
		r, err := applyJMESPath(query, result, rawOutput)
		if err != nil {
			fmt.Println(err)
		}
		result = r
	}
	fmt.Println(string(result))
}

func createAddCmd() *cobra.Command {
	opts := &AddOptions{}

	cmd := &cobra.Command{
		Use:           "add",
		Short:         "",
		Args:          cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := AddSilence(opts.EndPoint, os.Stdin)
			if err != nil {
				return err
			}
			printResult(result, opts.Query, opts.RawOutput)
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.EndPoint, "endpoint", "e", "", "API endpoint URL (ex: https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/)")
	cmd.MarkFlagRequired("endpoint")
	cmd.Flags().StringVarP(&opts.Query, "query", "q", "", "JMESPath query (ex: 'silenceID')")
	cmd.Flags().BoolVarP(&opts.RawOutput, "raw-output", "r", false, "print string as raw output")

	return cmd
}

func createDeleteCommand() *cobra.Command {
	opts := &DeleteOptions{}

	cmd := &cobra.Command{
		Use:           "delete",
		Short:         "",
		Args:          cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := DeleteSilence(opts.EndPoint, opts.SilenceId)
			if err != nil {
				return err
			}
			printResult(result, opts.Query, opts.RawOutput)
			return nil
		},
	}
	cmd.Flags().StringVarP(&opts.SilenceId, "silenceid", "s", "", "silence id (ex: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)")
	cmd.Flags().StringVarP(&opts.EndPoint, "endpoint", "e", "", "API endpoint URL (ex: https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/)")
	cmd.MarkFlagRequired("endpoint")
	cmd.Flags().StringVarP(&opts.Query, "query", "q", "", "JMESPath query (ex: 'silenceID')")
	cmd.Flags().BoolVarP(&opts.RawOutput, "raw-output", "r", false, "print string as raw output")

	return cmd
}

func createListCmd() *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:           "list",
		Short:         "",
		Args:          cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := ListSilences(opts.EndPoint)
			if err != nil {
				return err
			}
			printResult(result, opts.Query, opts.RawOutput)
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.EndPoint, "endpoint", "e", "", "API endpoint URL (ex: https://aps-workspaces.ap-northeast-1.amazonaws.com/workspaces/ws-xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx/)")
	cmd.MarkFlagRequired("endpoint")
	cmd.Flags().StringVarP(&opts.Query, "query", "q", "", "JMESPath query (ex: 'silenceID')")
	cmd.Flags().BoolVarP(&opts.RawOutput, "raw-output", "r", false, "print string as raw output")

	return cmd
}

func createListWorkspcesCmd() *cobra.Command {
	opts := &ListWorkspacesOptions{}

	cmd := &cobra.Command{
		Use:           "list-workspaces",
		Short:         "",
		Args:          cobra.MatchAll(cobra.NoArgs, cobra.OnlyValidArgs),
		SilenceUsage:  true,
		SilenceErrors: false,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := ListWorkspaces()
			if err != nil {
				return err
			}
			printResult(result, opts.Query, opts.RawOutput)
			return nil
		},
	}

	cmd.Flags().StringVarP(&opts.Query, "query", "q", "", "JMESPath query (ex: 'silenceID')")
	cmd.Flags().BoolVarP(&opts.RawOutput, "raw-output", "r", false, "print string as raw output")

	return cmd
}
