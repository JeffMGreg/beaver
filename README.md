Beaver
======
![alt tag](https://raw.github.com/JeffMGreg/beaver/master/beaver.jpg)

## Description
A simple GoLang logging package that supports levels and colors when writing to the terminal.  A project that I'm writing to help me learn Go.  Probably shouldn't use in production.

## Usage
```go
package main
import "beaver"

func main(){
  // Let's make a new logger
  logger := beaver.NewLogger(nil, "main")
  
  // I want color on my output
  logger.EnableColor()
  
  // Logging with different levels
  logger.Debug("This is a debug message")
  logger.Info("and this is an info message")
  logger.Warn("can you guess what level this is?")
  logger.Error("I think you get the picture")
  
  // Lets set the level higher
  logger.SetLevel(beaver.WARNING)
  
  logger.Info("Won't be printed")
  logger.Warn("Only printing warnings and above")
  
  // I'm tired of these colors
  logger.DisableColor()
  
  logger.Fatal("Oh no... something Fatal happened!")
}
```
