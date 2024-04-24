import http from 'k6/http';

var payloads =  [
  {
    id: 'test0',
    data: 'aaa',
  },
  {
    id: 'test1',
    data: 'bbb',
  },
  {
    id: 'test2',
    data: 'ccc',
  },
  {
    id: 'test3',
    data: 'ddd',
  },
  {
    id: 'test4',
    data: 'eee',
  },
  {
    id: 'test5',
    data: 'fff',
  }]

const url = 'http://localhost:8080';


export default function () {

  init();
  list();
  del();
  update();
}


function init() {

  payloads.forEach(function(obj) {
    const params = {
      headers: {
        'Content-Type': 'application/json',
      },
    };
  
    http.post(url, JSON.stringify(obj), params);
  })

}

function create() {
  const url = 'http://localhost:8080';
  const payload = JSON.stringify({
    id: 'test',
    data: 'bbb',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  http.post(url, payload, params);
}

function del() {

  http.del(url+"/delete/test1");
}

function update() {
  const payload = JSON.stringify({
    id: 'test1',
    data: 'aawef',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  http.post(url+"/update", payload, params);
}

function read() {

  http.get(url+"/read/test0");
  http.get(url+"/read/test1");
  http.get(url+"/read/test2");
}

function list() {
  http.get(url+"/list")
}