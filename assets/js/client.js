

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
            console.log("Соединение установлено.");
        };

        socket.onclose = function(event) {
            if (event.wasClean) {
                console.log('Соединение закрыто чисто');
            } else {
                console.log('Обрыв соединения'); // например, "убит" процесс сервера
            }
            console.log('Код: ' + event.code + ' причина: ' + event.reason);
        };

        socket.onmessage = (event) => {
            console.log("Получены данные " + event.data);
            this.$data.messages.push({
                author: 'Client 1',
                text: event.data
            });
        };

        socket.onerror = function(error) {
            console.log("Ошибка " + error.message);
        };

        this.$data.socket = socket;
    }
});