# Cobra入门

本文档基于[Cobra-UserGuide](https://github.com/spf13/cobra/blob/main/user_guide.md#positional-and-custom-arguments)

测试项目：[myapp](./test/myapp)

## 一、项目初始化

初始化一个MyApp项目，以供入门学习：

```bash
$ mkdir myapp
$ cd myapp
$ go mod init MyApp
$ cobra-cli init --author jiadisu@fudan.edu.cn --license apache
```

初始化完成，运行：

```bash
$ go run main.go
# 会显示一个简单的说明，这里没复制
```

## 二、编写root.go

`rootCmd`是没有指定任何子命令时执行的命令，指定子命令的示例：

```bash
$ go run main.go serve	# 执行子命令serve
$ go run main.go		# 这个就是没有指定任何子命令，会执行root
```

### 1、修改rootCmd

简单修改一下root的short和long字段，并把`Run`的注释解开，执行一句输出：

```go
// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "MyApp",
	Short: "MyApp is a learning program for getting started with Cobra.",
	Long: `Based on cobra-UserGuide, follow along and learn to create instructions, 
	set flag, parameter validator and other operations 
	to complete the basic Cobra operation introduction.`,
	Run: func(cmd *cobra.Command, args []string) {		// 指令的执行函数
		fmt.Println(cmd.Short)
	},
}
```

此时测试一下，执行如下命令：

```bash
$ go run main.go 
MyApp is a learning program for getting started with Cobra.
```

显示出Short中的那句话就算成功，说明不指定任何子命令时，会调用root的Run。

（另外废话一句，`fmt.Println(cmd.Short)`需要导入"fmt"模块，如果vscode配置好环境的话，这里不需要手动导入，需要时会自动导入）

### 2、添加flags

flags可以指令在执行时需要指定的参数，例如设置了一个flag是`str`，那么在指令执行时就可以有如下调用：

```bash
$ cmd --str strValue
```

该参数`strValue`就可以在执行Run时使用，从而达到不同的指令执行效果。

添加flags需要定义一个全局的变量，用来保存传入的flags的值，例如上面的`str`标签，就需要设置一个变量`var strV`与之对应，这样键入`cmd --str strValue`时，`strValue`就会被传到`var strV`变量中，以供编程使用。

flags分为Persistent Flags和Local Flags，即前者是全局的，所有指令（父子关系）可用，后者是局部的，只有自己可用。

举个实例，添加指令flags：

```bash
$ cobra-cli add flags --author jiadisu@fudan.edu.cn --license apache
```

在flags.go中改造：

```go
// 添加两个全局变量，以供flags使用
var Name string
var Save bool

// 在init函数中，初始化cmd的属性（包括flags）
func init() {
    // root为其parent
	rootCmd.AddCommand(flagsCmd)
	// name flag
    // StringVar(pointer, name, default value, usage)
	flagsCmd.Flags().StringVar(&Name, "name", "default name", "my name")
	// save flag
	flagsCmd.Flags().BoolVar(&Save, "save", false, "for save")
}
```

修改Run，以便看到flag设置的效果：

```go
Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("flags called")
    fmt.Println("cmd.Name = " + Name)
    if Save {
        fmt.Println("cmd.Save is true")
    }
},
```

进行如下指令测试：

```bash
$ go run main.go flags 							# 不设置name和save
flags called
cmd.Name = default name
$ go run main.go flags --name TestName			# 设置name为TestName
flags called
cmd.Name = TestName
$ go run main.go flags --name TestName --save	#设置name为TestName并设置save
flags called
cmd.Name = TestName
cmd.Save is true
```

### 3、配置config

config即应用的配置，由[viper](https://github.com/spf13/viper)项目管理，其实就是有点像一个功能丰富过的哈希表，记录每个配置的key与其对应的value，使用的函数也是比较自解释的。

例如：

```go
// 先设置flag，然后将其与config绑定（bind）
func init() {
  rootCmd.PersistentFlags().StringVar(&author, "author", "YOUR NAME", "Author name for copyright attribution")
  viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
}
```

### 4、flags的一些属性

#### （1）必须的flags

如果想要将某个flag指定为必须的，那么可通过`MarkFlagRequired()`实现，在原来的项目上改造如下：

```go
func init() {
	rootCmd.AddCommand(flagsCmd)
	// name flag
	flagsCmd.Flags().StringVar(&Name, "name", "default name", "my name")
	flagsCmd.MarkFlagRequired("name")		// 添加required
	// save flag
	flagsCmd.Flags().BoolVar(&Save, "save", false, "for save")
}
```

此时按如下方式执行：

```bash
$ go run main.go flags							# 不指定name，会产生如下报错
Error: required flag(s) "name" not set
Usage:
  MyApp flags [flags]

Flags:
  -h, --help          help for flags
      --name string   my name (default "default name")
      --save          for save

exit status 1
$ go run main.go flags --name TestName				# 指定name，执行成功
flags called
cmd.Name = TestName
```

#### （2）相互关联的flags

有时可能有两个flag可能**需要同时指定**，只指定一个无法执行（如登陆时username和password必须同时指定），可通过`MarkFlagsRequiredTogether()`将两个flag绑定起来，在原来的项目上改造如下：

```go
var Password string

func init() {
	rootCmd.AddCommand(flagsCmd)
	// name flag
	flagsCmd.Flags().StringVar(&Name, "name", "default name", "my name")
	flagsCmd.MarkFlagRequired("name")
	// password flag
	flagsCmd.Flags().StringVar(&Password, "password", "", "my password")
	flagsCmd.MarkFlagsRequiredTogether("name", "password")			// password与username绑定起来
	// save flag
	flagsCmd.Flags().BoolVar(&Save, "save", false, "for save")
}
```

此时执行如下测试：

```bash
$ go run main.go flags --name TestName							# 未指定password
Error: if any flags in the group [name password] are set they must all be set; missing [password]
Usage:
  MyApp flags [flags]

Flags:
  -h, --help              help for flags
      --name string       my name (default "default name")
      --password string   my password
      --save              for save

exit status 1
$ go run main.go flags --name TestName --password TestPassword		# 指定password，成功执行
flags called
cmd.Name = TestName
```

有时可能有两个flag可能**不能同时指定**，如格式的选项`json`和`yaml`不能同时指定，可通过`MarkFlagsMutuallyExclusive()`将两个flag绑定起来，在原来的项目上改造如下：

```go
var Json bool
var Yaml bool

var flagsCmd = &cobra.Command{
	Use:   "flags",
	Short: "Test flags",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("flags called")
		fmt.Println("cmd.Name = " + Name)
		if Save {
			fmt.Println("cmd.Save is true")
		}
		if Json {
			fmt.Println("cmd.format is JSON")
		} else {
			fmt.Println("cmd.format is YAML")
		}
	},
}

func init() {
	rootCmd.AddCommand(flagsCmd)
	// name flag
	flagsCmd.Flags().StringVar(&Name, "name", "default name", "my name")
	flagsCmd.MarkFlagRequired("name")
	// password flag
	flagsCmd.Flags().StringVar(&Password, "password", "", "my password")
	flagsCmd.MarkFlagsRequiredTogether("name", "password")
	// save flag
	flagsCmd.Flags().BoolVar(&Save, "save", false, "for save")
	// json and yaml flag
	flagsCmd.Flags().BoolVar(&Json, "json", false, "format json")
	flagsCmd.Flags().BoolVar(&Yaml, "yaml", false, "format yaml")
	flagsCmd.MarkFlagsMutuallyExclusive("json", "yaml")		// json与yaml不能同时指定
}
```

执行如下测试：

```bash
$ go run main.go flags --name TestName --password TestPassword		# 不指定json与yaml
flags called
cmd.Name = TestName
cmd.format is YAML
$ go run main.go flags --name TestName --password TestPassword --json	# json
flags called
cmd.Name = TestName
cmd.format is JSON
$ go run main.go flags --name TestName --password TestPassword --json --yaml	# 指定json与yaml，如下报错
Error: if any flags in the group [json yaml] are set none of the others can be; [json yaml] were all set
Usage:
  MyApp flags [flags]

Flags:
  -h, --help              help for flags
      --json              format json
      --name string       my name (default "default name")
      --password string   my password
      --save              for save
      --yaml              format yaml

exit status 1
```

### 5、Arguments参数

#### （1）arguments是什么

Arguments就是指令中**除了flags外传入的字段**，会被保存在一个`args[]`数组中，供Run使用。

做个小测试，新建指令args并修改Run，令其输出`args[]`内容：

```go
Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("args called")
    fmt.Printf("args: %v\n", args)
},
```

执行如下指令：

```bash
$ go run main.go args argv1 argv2		# 包含了两个argument: argv1和argv2
args called
args: [argv1 argv2]
```

#### （2）为arguments添加校验

cobra提供有如下几种args校验器：

- 参数个数
    - `NoArgs` - 限制没有参数
    - `ArbitraryArgs` - 任何参数都可以
    - `MinimumNArgs(int)` - 指定最少的参数个数
    - `MaximumNArgs(int)` - 指定最多的参数个数
    - `ExactArgs(int)` - 指定明确的参数个数
    - `RangeArgs(min, max)` - 指定参数个数范围
- 内容:
    - `OnlyValidArgs` - report an error if there are any positional args not specified in the `ValidArgs` field of `Command`, which can optionally be set to a list of valid values for positional args.

If `Args` is undefined or `nil`, it defaults to `ArbitraryArgs`.

另外，`MatchAll()`可以将多个校验器结合起来。

为cmd添加校验的方式：指定`cmd.Args`字段：

```go
var argsCmd = &cobra.Command{
	Use:   "args",
	Short: "Test arguments",
	Args: cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),		// 校验
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("args called")
		fmt.Printf("args: %v\n", args)
	},
}
```

测试：

```bash
$ go run main.go args						# 参数个数不对
Error: accepts 2 arg(s), received 0
Usage:
  MyApp args [flags]

Flags:
  -h, --help   help for args

exit status 1
$ go run main.go args argv1 argv2			# 校验通过
args called
args: [argv1 argv2]
```

#### （3）自定义校验

除了cobra提供的一系列校验外，可以自定义校验，将Args设置为如下声明的函数：

```go
func(cmd *cobra.Command, args []string) error {}
```

例如，创建一个blue指令：

```bash
$ cobra-cli add blue --author jiadisu@fudan.edu.cn --license apache
```

修改其变量字段如下：

```go
// blueCmd represents the blue command
var blueCmd = &cobra.Command{
	Use:   "blue",
	Short: "An example to test args",
	Args: func(cmd *cobra.Command, args []string) error {
		err := cobra.MinimumNArgs(1)(cmd, args)		// 校验个数
		if err != nil {
			return err // 不满足min要求
		}
		if args[0] == "blue" {
			return nil // 通过校验
		}
		return fmt.Errorf("invalid color specified: %s", args[0])
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("blue called")
		fmt.Printf("args: %v\n", args)
	},
}
```

上面的args只有在个数大于1且第一个内容为blue时才会通过校验：

```bash
$ go run main.go blue						# 缺少参数
Error: requires at least 1 arg(s), only received 0
Usage:
  MyApp blue [flags]

Flags:
  -h, --help   help for blue

exit status 1
$ go run main.go blue argv					# 不是blue
Error: invalid color specified: argv
Usage:
  MyApp blue [flags]

Flags:
  -h, --help   help for blue

exit status 1
$ go run main.go blue blue lalala			# 通过校验
blue called
args: [blue lalala]
```

## 三、实例

sample.go文件内容：

```go
package main

import (
    "fmt"
    "strings"

    "github.com/spf13/cobra"
)

func main() {
    var echoTimes int

    var cmdPrint = &cobra.Command{							// print
        Use:   "print [string to print]",
        Short: "Print anything to the screen",
        Long: `print is for printing anything back to the screen.
        For many years people have printed back to the screen.`,
        Args: cobra.MinimumNArgs(1),							// 最少1个
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Print: " + strings.Join(args, " "))		// " "分隔，逐个输出
        },
    }

    var cmdEcho = &cobra.Command{								// echo
        Use:   "echo [string to echo]",
        Short: "Echo anything to the screen",
        Long: `echo is for echoing anything back.
        Echo works a lot like print, except it has a child command.`,
        Args: cobra.MinimumNArgs(1),							// 最少1个
        Run: func(cmd *cobra.Command, args []string) {
            fmt.Println("Echo: " + strings.Join(args, " "))		// " "分隔，逐个输出
        },
    }

    var cmdTimes = &cobra.Command{							// times
        Use:   "times [string to echo]",
        Short: "Echo anything to the screen more times",
        Long: `echo things multiple times back to the user by providing
        a count and a string.`,
        Args: cobra.MinimumNArgs(1),							// 最少1个
        Run: func(cmd *cobra.Command, args []string) {
            for i := 0; i < echoTimes; i++ {
                fmt.Println("Echo: " + strings.Join(args, " "))		// 输出times次
            }
        },
    }

    // times的flag -t --times
    cmdTimes.Flags().IntVarP(&echoTimes, "times", "t", 1, "times to echo the input")

    var rootCmd = &cobra.Command{Use: "app"}
    rootCmd.AddCommand(cmdPrint, cmdEcho)
    cmdEcho.AddCommand(cmdTimes)					// times为echo的子命令
    rootCmd.Execute()
}
```

简单测试：

```bash
$ go run sample.go print
Error: requires at least 1 arg(s), only received 0
Usage:
  app print [string to print] [flags]

Flags:
  -h, --help   help for print
$ go run sample.go print hello world
Print: hello world
$ go run sample.go echo hello world
Echo: hello world
$ go run sample.go times hello world, -t 3
Error: unknown command "times" for "app"
Run 'app --help' for usage.
$ go run sample.go -h
Usage:
  app [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  echo        Echo anything to the screen
  help        Help about any command
  print       Print anything to the screen

Flags:
  -h, --help   help for app

Use "app [command] --help" for more information about a command.
$ go run sample.go echo -h
echo is for echoing anything back.
Echo works a lot like print, except it has a child command.

Usage:
  app echo [string to echo] [flags]
  app echo [command]

Available Commands:
  times       Echo anything to the screen more times

Flags:
  -h, --help   help for echo

Use "app echo [command] --help" for more information about a command.
$ go run sample.go echo times hello world -t 3
Echo: hello world
Echo: hello world
Echo: hello world
```

