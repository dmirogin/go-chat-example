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
        let socket = new WebSocket('ws://' + document.location.host + '/ws');

        socket.onopen = function() {
            console.log("Connection established.");
        };

        socket.onclose = function(event) {
            if (event.wasClean) {
                console.log('Clean close');
            } else {
                console.log('Break connection'); // for example, server stopped connection
            }
            console.log('Errore code: ' + event.code + ' reason: ' + event.reason);
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


        const eventSource = new EventSource("/sse");

        eventSource.onopen = function(e) {
            console.log("Соединение открыто");
        };

        eventSource.onerror = function(e) {
            if (this.readyState == EventSource.CONNECTING) {
                console.log("Соединение порвалось, пересоединяемся...");
            } else {
                console.log("Ошибка, состояние: " + this.readyState);
            }
        };

        eventSource.onmessage = function(e) {
            console.log("Пришли данные: " + e.data);
        };

        eventSource.addEventListener('time', function(e) {
            console.log('Пришёл ' + e.data );
        });

    }
});