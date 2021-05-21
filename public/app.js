new Vue({
    el: '#app',

    data: {
        ws: null, // websocket
        newMsg: '', // holds new msgs to be sent to the server
        chatContent: '', // running list of chat messages displayed
        email: null, // email address used for grabbing a gravatar
        username: null, // username
        joined: false // true if email/username have been filled
    },

    created: function() {
        var self = this;
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        this.ws.addEventListener('message', function(e) {
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                + '<img src="' + self.gravatarURL(msg.email) + '">' // avatar
                + msg.username
            + '</div>'
            + emojione.toImage(msg.message) + '<br/>'; // parse emojis
            
        var element = document.getElementById('chat-messages');
        element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        });
    },

    methods: {
        send: function() {
            if (this.newMsg != '') {
                this.ws.send(
                    JSON.stringify({
                        email: this.email,
                        username: this.username,
                        message: $('<p>').html(this.newMsg).text() // strip out html
                    }
                    ));
                    this.newMsg = ''; // reset newMsg
            }
        },

        join: function() {
            if (!this.email) {
                Materialize.toast('Email required!', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('Username required!', 2000);
                return
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },

        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});