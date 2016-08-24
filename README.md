# snapfs
_Currently a work in progress._

Snapfs lets you watch local paths on your filesystem for changes to its files and subdirectories. The current state of all watched paths, when summed, is called a ``Snapshot``. These snapshots are marshaled into JSON format that you can later ``Restore`` to snapfs. 

Snapfs relies on polling and its up to you to regularly call ``Update`` when you want the latest changes. All changes since the prior call to ``Update`` is returned as an ``Events`` struct.

You can specify which files/folders to ignore by using regular expressions with ``Ignore``. 

I wrote this to satisfy a particular need in another project of mine. I've yet to test this fully nor can I assure you that it is performant. Please don't rely on this for production or mission critical applications. I mostly wrote this for fun and excerise as I am still very new to Go.

**TODO LIST**
- Finishing implementing tests.
- Use worker pools and goroutines within ``Update`` -- could definitely use guidance on this
- Enhance documentation and this readme

I very much welcome all issues, advice, and especially pull requests. 