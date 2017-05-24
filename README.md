# realstream
Pulls term based data from Twitter's realtime streaming API.

Backend:
MongoDB, NSQ, go.

Frontend:
Bootstrap, Google's webkit, Javascript.


• We can distribute our MongoDB and NSQ nodes across many
physical machines which would mean our system is capable of gigantic
scale—whenever resources start running low, we can add new boxes to
cope with the demand.

• When we add other applications that need to query and read the results
from polls, we can be sure that our database services are highly available
and capable of delivering.

• We can spread our database across geographical expanses replicating
data for backup so we don't lose anything when disaster strikes.

• We can build a multi-node, fault tolerant NSQ environment, which means
when our twittervotes program learns of interesting tweets, there will
always be somewhere to send the data.

• We could write many more programs that generate votes from different
sources; the only requirement is that they know how to put messages
into NSQ.
