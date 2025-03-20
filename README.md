# TypeKit

TypeKit is a lightweight dependency injection framework for Go. It works by allowing the developer to register types using typekit.Register(), ensuring that each type is instantiated through its own constructor. Within these constructors, typekit.Resolve() takes care of fetching dependencies when they are needed, simplifying dependency management.

## 🔧 Features  

- **Auto-Inject Dependencies** – Call `typekit.Resolve()` inside constructors, and TypeKit handles the rest.  
- **No Manual Wiring** – No need to pass dependencies around manually.  
- **Simple API** – Register constructors, then resolve types anywhere.  
- **Error Handling** – Uses `func() (T, error)` constructors for safety.  

---

## 📦 Installation  

```sh
go get github.com/the-zucc/typekit
```

---

## 📝 Usage  

### 1️⃣ Define Dependencies  

```go
package main

import (
	"fmt"
	"github.com/the-zucc/typekit"
)

type A struct{}

func NewA() (A, error) {
	fmt.Println("Creating A")
	return A{}, nil
}

type B struct {
	A *A
}

func NewB() (B, error) {
	a := typekit.Resolve[A]() // Auto-resolve A
	fmt.Println("Creating B")
	return B{A: a}, nil
}

type C struct {
	A *A
	B *B
}

func NewC() (C, error) {
	a := typekit.Resolve[A]() // Auto-resolve A
	b := typekit.Resolve[B]() // Auto-resolve B
	fmt.Println("Creating C")
	return C{A: a, B: b}, nil
}
```

---

### 2️⃣ Register Everything in `main.go`  

```go
func main() {
	// Register dependencies
	typekit.Register(NewA)
	typekit.Register(NewB)
	typekit.Register(NewC)

	// Resolve final instance
	c := typekit.Resolve[C]()
	fmt.Printf("Resolved instance: %#v\n", c)
}
```

---

## 🔍 Why This Works  

- **No function arguments** – Dependencies are resolved inside constructors.  
- **Minimal setup** – Just register functions and resolve when needed.  
- **Works across packages** – No need to pass dependencies manually.  

---

## ⚠️ Error Handling  

If TypeKit panics, it's usually because a dependency **wasn’t registered** before calling `Resolve()`.  

---

## 📜 License  

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.  