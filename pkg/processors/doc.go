/*
Package processors contains functions to process items.

The processor takes a pointer to the item and to a waitgroup.
It is the responsibility of the processor to Add to the waitgroup and to call Done() when it is done processing the item.
Failing to do so will cause the program to behave in unpredictable ways.

The processor returns a boolean and an error.
The boolean is used to signal to the caller if it should exclude the kids of the item from the processing when an error is returned.

Example:

	// This processor will return an error if the item has no text
	// and it will exclude the kids of the item from the processing
	processor := func(item *gohn.Item, wg *sync.WaitGroup) (bool, error) {
		wg.Add(1)
		defer wg.Done()
		if item.Text == nil {
			return true, fmt.Errorf("item has no text")
		}
		return false, nil
	}

If you want the kids of the item to be processed, you can return false as the first return value.

Example:

	// This processor will return an error if the item has no text
	// and it will include the kids of the item in the processing
	processor := func(item *gohn.Item, wg *sync.WaitGroup) (bool, error) {
		wg.Add(1)
		defer wg.Done()
		if item.Text == nil {
			return false, fmt.Errorf("item has no text")
		}
		return false, nil
	}

If you need to further process the items, you can send the item to a channel.
Example:

	processor := func(item *gohn.Item, wg *sync.WaitGroup) (bool, error) {
		wg.Add(1)
		defer wg.Done()
		if item.Text == nil {
			return false, fmt.Errorf("item has no text")
		}
		ch <- item
		return false, nil
	}

*/

package processors
