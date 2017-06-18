# Realstream

Pulls term based data from Twitter's realtime streaming API.

Backend:
MongoDB, NSQ, go.

Frontend:
Bootstrap, Google Webkit, Javascript.

![alt text](https://raw.githubusercontent.com/rflorenc/realstream/master/res/realstream.png)



Example app.


• Distributing MongoDB and NSQ nodes across many containers
  containers would mean the app is capable of gigantic scale.

• Spread our database across geographical regions replicating
data for backup so we don't lose anything when disaster strikes.

• Build a multi-node/multi container, fault tolerant NSQ environment, which means
when twittervotes program learns of interesting votes, there will always be somewhere to send the data.
