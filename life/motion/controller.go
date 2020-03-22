package motion

import "github.com/fogleman/ease"

// A Controller controls the motion of a living thing.
type Controller struct {
	// Current and next intents.
	current, next Intent

	// The progress and duration of the current intent.
	progress, duration int

	// A multiplier that affects the speed of all motions.
	speed float64

	// The current velocity.
	dx, dy float64
}

// NewController creates a Controller.
func NewController() *Controller {
	return &Controller{
		duration: 12,
		speed:    1,
	}
}

// Intent returns the current Intent.
func (c *Controller) Intent() Intent { return c.current }

// Next returns true if the Controller's next Intent is set.
func (c *Controller) Next() bool { return c.next != Stay }

// Update is called once a tick to update the Controller.
func (c *Controller) Update(next Intent) {
	v := ease.OutQuint(float64(c.progress) / float64(c.duration))

	c.dy = 0
	c.dx = 0
	switch c.current {
	case StepUp:
		c.dy = -v * c.speed
	case StepDown:
		c.dy = v * c.speed
	case StepLeft:
		c.dx = -v * c.speed
	case StepRight:
		c.dx = v * c.speed
	}

	// Update intent.
	if c.next == Stay {
		c.next = next
	}

	// Update progress.
	c.progress++
	if c.progress == c.duration {
		c.progress = 0
		c.current = c.next
		c.next = Stay
	}
}

// Motion returns the current velocity of the controller.
func (c *Controller) Motion() (dx, dy float64) { return c.dx, c.dy }
