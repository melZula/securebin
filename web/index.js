var getParams = function () {
  var params = {};
  var query = location.search.substring(1);
  var vars = query.split('&');
  for (var i = 0; i < vars.length; i++) {
    var pair = vars[i].split('=');
    params[pair[0]] = decodeURIComponent(pair[1]);
  }
  return params;
};

var app = new Vue({
  el: '#app',
  data: {
    message: 'Hello world!',
    lifetime: 0,
    id: 0,
    img: '',
    pass: '',
    reqsts: [],
    show: false,
    auth: true,
  },

  updated: function () {
    if (!this.show) {
      $('select').formSelect();
    }
  },
  created: function () {
    this.parseArgs();
  },
  methods: {
    reverseMessage: function () {
      this.message = this.message.split('').reverse().join('');
    },
    sendText: function () {
      let payload = {
        text: this.message,
        time: parseInt(this.lifetime),
      };

      fetch('http://securebin.local/api/paste', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
        .then((response) => {
          if (response.ok) {
            return response.json();
          } else {
            return Promise.reject(response.status);
          }
        })
        .then((r) => {
          this.id = r.id;
          this.pass = r.password;
          alert('OK');
        })
        .catch((err) => {
          alert(err.status + ': ' + err.statusText);
        });
    },
    parseArgs: function () {
      let params = getParams();
      if (params.id > 0) {
        this.id = params.id;
        this.show = params ? true : false;
      }
    },
    getImage: function () {
      this.auth = false;
      let payload = {
        id: parseInt(this.id),
        password: this.pass,
      };
      fetch('http://securebin.local/api/data', {
        method: 'POST',
        body: JSON.stringify(payload),
      })
        .then((response) => {
          console.log(response);
          if (response.ok) {
            return response.json();
          } else {
            return Promise.reject(response.status);
          }
        })
        .then((json) => {
          this.img = 'data:image/png;base64,' + json.img.replace(/^"|"$/gm, '');
          this.reqsts = json.times;
        })
        .catch((err) => {
          console.log(err);
          alert(err.status + ': ' + err.statusText);
        });
    },
  },
});
