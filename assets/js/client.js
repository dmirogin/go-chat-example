

new Vue({
    el: '#root',
    data: {
        message: '',
        messages: [
            {
                author: 'client 1',
                text: 'Text'
            },
            {
                author: 'client 2',
                text: 'Text 2'
            }
        ],
        socket: ''
    },
    methods: {
        sendMessage() {
            console.log(this.$data.message);
            this.$data.socket.send(this.$data.message);
            this.$data.message = '';
        }
    },
    created() {
        let socket = new WebSocket("ws://localhost:3000/echo");

        socket.onopen = function() {
            console.log("Connection established.");
        };

        socket.onclose = function(event) {
            if (event.wasClean) {
                console.log('Clean close');
            } else {
                console.log('Break connection'); // for example, server stopped connection
            }
            console.log('Код: ' + event.code + ' причина: ' + event.reason);
        };

        socket.onmessage = (event) => {
            console.log("Received data", event.data);
            let data = JSON.parse(event.data);
            this.$data.messages.push({
                author: 'Client ' + data.author,
                text: data.text
            });
        };

        socket.onerror = function(error) {
            console.log("Error " + error.message);
        };

        this.$data.socket = socket;
    }
});