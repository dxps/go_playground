## Go Practice

This is a practice field for various standard Golang topics and features.

<br/>

### Math > Quadratic Equation

Implement the function `findRoots` to find the roots of the quadratic equation: ax2 + bx + c = 0. If the equation has only one solution, the function should return that solution as both results. The equation will always have at least one solution.

The roots of the quadratic equation can be found with the following formula:

![A quadratic equation.](./equation.png)

For example, the roots of the equation 2x<sup>2</sup> + 10x + 8 = 0 are -1 and -4.

<br/>

### Slices > Append when space

The _"append when space"_ example shows how slices are working internally, meaning that if there is enough capacity `append` function is just appending the data into the internal array, and it doesn't create a new array to accomodate the space (capacity) requirement.

That's the reason, at the end, both `x` and `y` are slices pointing to the same underlying array, as confirmed by the execution output:
```shell
$ go run append_when_space.go
 x = [g o]    internals:  Len=2  Cap=3  Data=824634188112
 y = [g o t]  internals:  Len=3  Cap=3  Data=824634188112
$ 
```
