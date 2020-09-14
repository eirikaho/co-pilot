package cmd

import (
	"co-pilot/pkg/logger"
	"co-pilot/pkg/upgrade"
	"fmt"
	"github.com/perottobc/mvn-pom-mutator/pkg/pom"
	"github.com/spf13/cobra"
)

var groupId string
var artifactId string

var upgradeCmd = &cobra.Command{
	Use:   "upgrade [OPTIONS]",
	Short: "Upgrade options",
	Long:  `Perform upgrade on existing projects`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := EnableDebug(cmd, args); err != nil {
			log.Fatalln(err)
		}
		populatePomFiles()
	},
}

var upgradeSpringBootCmd = &cobra.Command{
	Use:   "spring-boot",
	Short: "Upgrade spring-boot to the latest version",
	Long:  `Upgrade spring-boot to the latest version`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading spring-boot for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}

			if err = upgrade.SpringBoot(model); err != nil {
				log.Warnln(err)
				continue
			}

			if !cArgs.DryRun {
				if err = write(pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

var upgradeDependencyCmd = &cobra.Command{
	Use:   "dependency",
	Short: "Upgrade a specific dependency on a project",
	Long:  `Upgrade a specific dependency on a project`,
	Run: func(cmd *cobra.Command, args []string) {
		if groupId == "" || artifactId == "" {
			log.Fatal("--groupId (-g) and --artifactId (-a) must be set")
		}
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading dependency %s:%s for pom file %s", groupId, artifactId, pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}

			if err := upgrade.Dependency(model, groupId, artifactId); err != nil {
				log.Warnln(err)
				continue
			}

			if !cArgs.DryRun {
				if err = write(pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

var upgrade2partyDependenciesCmd = &cobra.Command{
	Use:   "2party",
	Short: "Upgrade all 2party dependencies to project",
	Long:  `Upgrade all 2party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading 2party for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}
			if err = upgrade.Dependencies(model, true); err != nil {
				log.Warnln(err)
				continue
			}

			if !cArgs.DryRun {
				if err = write(pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

var upgrade3partyDependenciesCmd = &cobra.Command{
	Use:   "3party",
	Short: "Upgrade all 3party dependencies to project",
	Long:  `Upgrade all 3party dependencies to project`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading 3party for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}

			if err = upgrade.Dependencies(model, false); err != nil {
				log.Warnln(err)
				continue
			}

			if !cArgs.DryRun {
				if err = write(pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

var upgradeKotlinCmd = &cobra.Command{
	Use:   "kotlin",
	Short: "Upgrade kotlin version in project",
	Long:  `Upgrade kotlin version in project`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading kotlin for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}

			if err = upgrade.Kotlin(model); err != nil {
				log.Warnln(err)
				continue
			}

			if !cArgs.DryRun {
				if err = write(pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

var upgradePluginsCmd = &cobra.Command{
	Use:   "plugins",
	Short: "Upgrade all plugins found in project",
	Long:  `Upgrade all plugins found in project`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading plugins for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}
			if err = upgrade.Plugin(model); err != nil {
				log.Warnln(err)
				continue
			}

			if !cArgs.DryRun {
				if err = write(pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

var upgradeAllCmd = &cobra.Command{
	Use:   "all",
	Short: "Upgrade everything in project",
	Long:  `Upgrade everything in project`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, pomFile := range cArgs.PomFiles {
			log.Info(logger.White(fmt.Sprintf("upgrading all for pom file %s", pomFile)))
			model, err := pom.GetModelFrom(pomFile)
			if err != nil {
				log.Warnln(err)
				continue
			}
			upgrade.All(model)

			if !cArgs.DryRun {
				if err := write(pomFile, model); err != nil {
					log.Warnln(err)
				}
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(upgradeCmd)
	upgradeCmd.AddCommand(upgradeDependencyCmd)
	upgradeCmd.AddCommand(upgrade2partyDependenciesCmd)
	upgradeCmd.AddCommand(upgrade3partyDependenciesCmd)
	upgradeCmd.AddCommand(upgradeSpringBootCmd)
	upgradeCmd.AddCommand(upgradeKotlinCmd)
	upgradeCmd.AddCommand(upgradePluginsCmd)
	upgradeCmd.AddCommand(upgradeAllCmd)

	upgradeDependencyCmd.PersistentFlags().StringVarP(&groupId, "groupId", "g", "", "GroupId for upgrade")
	upgradeDependencyCmd.PersistentFlags().StringVarP(&artifactId, "artifactId", "a", "", "ArtifactId for upgrade")

	upgradeCmd.PersistentFlags().BoolVarP(&cArgs.Recursive, "recursive", "r", false, "turn on recursive mode")
	upgradeCmd.PersistentFlags().StringVar(&cArgs.TargetDirectory, "target", ".", "Optional target directory")
	upgradeCmd.PersistentFlags().BoolVar(&cArgs.Overwrite, "overwrite", true, "Overwrite pom.xml file")
	upgradeCmd.PersistentFlags().BoolVar(&cArgs.DryRun, "dry-run", false, "dry run does not write to pom.xml")
}
