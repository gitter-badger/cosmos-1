use std::io::{Read, Write};
use std::net::{TcpListener, TcpStream};
use std::thread;
use std::env;

fn call() {
    // host
    let _ = match env::var("HOST") {
        Ok(host) => host,
        Err(e) => { println!("{}", e); return; }
    };

    // token
    let _ = match env::var("TOKEN") {
        Ok(token) => token,
        Err(e) => { println!("{}", e); return; }
    };

    // try to connect to curiosity

    // smtp notification
}

fn handle_client(s: TcpStream) {
    let mut stream = match s.try_clone() {
        Ok(stream) => stream,
        Err(e) => { println!("{}", e); return; }
    };

    const BUFFER_SIZE: usize = 1024;
    let mut buffer: [u8; BUFFER_SIZE] = [0; BUFFER_SIZE];
    let mut raw: Vec<u8> = Vec::new();
    loop {
        let len = match stream.read(&mut buffer) {
            Ok(len) => len,
            Err(_) => { break; }
        };
        for i in 0..len { raw.push(buffer[i]); }
        if len < BUFFER_SIZE { break; }
    }

    let request = match String::from_utf8(raw) {
        Ok(request) => request,
        Err(e) => { println!("{}", e); return; }
    };

    println!("{}", request);

    let response = "HTTP/1.1 200 OK\r\nContent-Type: application/json; charset=UTF-8\r\nConnection: close\r\n\r\n";
    let _ = stream.write(response.as_bytes());
    
    call();
}

fn main() {
    let listener = TcpListener::bind("0.0.0.0:8888").unwrap();
    for stream in listener.incoming() {
        match stream {
            Ok(stream) => {
                thread::spawn(move|| {
                    handle_client(stream)
                });
            }
            Err(e) => { println!("{}", e); }
        }
    }
    drop(listener);
}
