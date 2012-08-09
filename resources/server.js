var webRTC = require('webrtc.io').listen(app);
var colors = {};

function selectRoom(socket) {
  for (var room in servers) {
    console.log('***' + room);
    if (io.sockets.clients(room).length < 4) {
      socket.emit('send', room);
    }
    console.log(io.sockets.clients('' + room));
  }
}

webRTC.rtc.on('connection', function(rtc) {
  rtc.on('send_answer', function() {
  });
  rtc.on('disconnect', function() {
  });
});

webRTC.sockets.on('connection', function(socket) {
  console.log("connection received");

  colors[socket.id] = Math.floor(Math.random()* 0xFFFFFF)
  socket.on('chat msg', function(msg) {
    console.log("chat received");
    
    socket.broadcast.emit('receive chat msg', {
      msg: msg,
      color: colors[socket.id]
    });
  });
});
