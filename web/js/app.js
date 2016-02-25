function getNotes() {
  var xmlhttp = new XMLHttpRequest();
  var url = '/api/notes';

  xmlhttp.onreadystatechange = function() {
      if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
	  var response = JSON.parse(xmlhttp.responseText);
          document.getElementById('d01').innerHTML += response + '<br>';
      }
  };
  xmlhttp.open('GET', url, true);
  xmlhttp.send();
}

function admin(token) {
  var xmlhttp = new XMLHttpRequest();
  var url = '/api/admin';

  xmlhttp.onreadystatechange = function() {
      if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
	  var response = JSON.parse(xmlhttp.responseText);
          document.getElementById('d01').innerHTML += 'response: ' + response.text + '<br>';
      }
  };

  xmlhttp.open('GET', url, true);
  xmlhttp.setRequestHeader('Authorization', 'Bearer '+token.token)
  xmlhttp.send();
}

function login() {
  var xmlhttp = new XMLHttpRequest();
  var url = '/api/login';

  xmlhttp.onreadystatechange = function() {
      if (xmlhttp.readyState == 4 && xmlhttp.status == 200) {
	  var response = JSON.parse(xmlhttp.responseText);
          document.getElementById('d01').innerHTML += response + '<br>';

          admin(response);
      }
  };
  xmlhttp.open('POST', url, true);

  var formData = {
      'username' : document.getElementById('username').value,
      'password' : document.getElementById('password').value
  };
  var data = JSON.stringify(formData);

  xmlhttp.send(data);
}

function main() {
  var source = new EventSource('/api/notes/events');
  source.onmessage = function(e) {
    document.getElementById('d01').innerHTML += e.data + '<br>';
  };
}
