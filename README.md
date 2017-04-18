# Alice [![Build Status](https://travis-ci.org/magic003/alice.png?branch=master)](https://travis-ci.org/magic003/alice)
Alice is an additive dependency injection container for Golang.

## Philosophy

Design philosophy behind Alice:
* The application components should **not** be aware of the existence of a DI container.
* Use static Go files to define the object graph.
* Developer has the freedom to choose the way to initialize objects.

## Install

```
$ go get github.com/magic003/alice
```

## Usage

Alice is inspired by the design of [Spring JavaConfig](http://docs.spring.io/spring-javaconfig/docs/1.0.0.M4/reference/html/). 

It usually takes 3 steps to use Alice.

### Define modules

The instances to be managed by the container are defined in modules. There could be multiple modules organized by the functionality of the instances. Modules are usually placed in a separate package.

A typical module looks like this:

```go
type ExampleModule struct {
    alice.BaseModule
    Foo Foo `alice:""`
    Bar Bar `alice:"Bar"`
    Baz Baz
}

func (m *ExampleModule) InstanceX() X {
    return X{m.Foo}
}

func (m *ExampleModule) InstanceY() Y {
    return Y{m.Baz}
}
```

A module struct must embed the `alice.BaseModule` struct. It allows 3 types of fields:
* Field tagged by `alice:""`. It will be associated with the same or assignable type of instance defined in other modules.
* Field tagged by `alice:"Bar"`. It will be associated with the instance named `Bar` defined in other modules.
* Field without `alice` tag. It will **not** be associated with any instance defined in other modules. It is expected to be provided when initializing the module. It is not managed by the container and could not be retrieved.

It is also common that no field is defined in a module struct.

Any public method of the module struct defines one instance to be intialized and maintained by the container. It is required to use a pointer receiver. The method name will be used as the instance name. The return type will be used as the instance type. Inside the method, it could use any field of the module struct to create new instances.

### Create container

During the bootstrap of the application, create a container by providing instances of modules.

```go
m1 := &ExampleModule1{}
m2 := &ExampleModule2{...}
container := alice.CreateContainer(m1, m2)
```

It will panic if any module is invalid.

### Retreive instances

The container provides 2 ways to retrieve instances: by name and by type.

```go
instanceX := container.InstanceByName("InstanceX")

instanceY := container.Instance(reflect.TypeOf((Y)(nil)))
```

It will panic either if no instance is found or if multiple matched types are found.

## Example

A dummy [example](https://github.com/magic003/alice/tree/master/example) using Alice.
