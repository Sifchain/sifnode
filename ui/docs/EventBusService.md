![Untitled Diagram (5)](https://user-images.githubusercontent.com/1256409/110735491-ce603f80-8279-11eb-8254-6a31e09e47eb.png)


This is an example of a decoupled notifications service based on an event emitter pattern.

* `EventBusService` is shared with the view
* Core dispatches events on the EventBusService
* View listens to those events and creates the appropriate view notification. 

Why is this a good idea?
* Core knows nothing about the messaging that the view uses to display the message. This will make things like localization easier.
* Google Analytics which is an implementation concern can get realtime information on what is happening within core by simply adding a listener to the NotificationsService. 
* Other implementation concerns such as logging and error reporting can listen to events driven from the EventBusService and marshall the appropriate response to their respective targets
* Events are typed so the event payload shape is matched to the event type
* Cohesion occurs through the shared definition of events which in this example happens [here](https://github.com/Sifchain/sifnode/pull/892/files#diff-1181e517a11ffd75848b4d3e55ccdaf88bf27ec325aa9c9ec5218d472e9d92e4R7) we might want to move them elsewhere.


Possible other things we could do
- [x] We could rename `NotificationsService` to be something more general such as `EventBusService` to be clearer 
- [x] We could rename `notify` to `dispatch` to be clearer.
- [ ] After we have renamed the structures within core this might look more like: `services.bus.dispatch(myevent)`

Possible other things we might be careful of doing
* We must not define events within core that are view specific events. Core should not know or care about events that are view specific. Instead we could trigger those events directly in code (eg `ga('send', 'pageview');`) but also avoid those events unless they are 100% necessary. 

Credits: @ryardley