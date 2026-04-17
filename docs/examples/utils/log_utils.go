package utils

import (
	"log"

	"github.com/seyallius/gosaidno/v2/aspect"
)

func LogBefore(c *aspect.Context, priority int, message string) {
	log.Printf("🟢 [BEFORE] %s - Priority: %d [%s]", c.FunctionName, priority, message)
}

func LogAfter(c *aspect.Context, priority int, message string) {
	log.Printf("🔵 [AFTER] %s - Priority: %d [%s]", c.FunctionName, priority, message)
}

func LogAround(c *aspect.Context, priority int, message string) {
	log.Printf("🟠 [AROUND] %s - Priority: %d [%s]", c.FunctionName, priority, message)
}

func LogAfterReturning(c *aspect.Context, priority int, message string) {
	log.Printf("🟣 [AFTER_RETURNING] %s - Priority: %d [%s]", c.FunctionName, priority, message)
}

func LogAfterThrowing(c *aspect.Context, priority int, message string) {
	log.Printf("🔴 [AFTER_THROWING] %s - Priority: %d [%s]", c.FunctionName, priority, message)
}
