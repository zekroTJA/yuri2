/** @format */

import { Injectable } from '@angular/core';
import { environment } from 'src/environments/environment';
import { EventEmitter } from 'events';
import { WSCommand, WSEvent } from './ws.static';
import { Message, EndEvent, WSErrorEvent } from './ws.models';
import { getCookieValue } from 'src/util/util.cookies';
import { ToastService } from 'src/app/components/toast/toast.service';

@Injectable({
  providedIn: 'root',
})
export class WSService extends EventEmitter {
  private rootURL: string;
  private ws: WebSocket;
  private _isInitialized: boolean;

  private getRootURLFromWindow(): string {
    const loc = window.location;
    return `${loc.href.startsWith('http://') ? 'ws' : 'wss'}://${loc.host}/ws`;
  }

  constructor(private toasts: ToastService) {
    super();

    console.log('OPENING WEB SOCKET CONNECTION');

    this.rootURL = environment.production
      ? this.getRootURLFromWindow()
      : 'ws://localhost:8080/ws';

    this.ws = new WebSocket(this.rootURL);

    this.ws.onmessage = (ev: MessageEvent): any => this.onMessage(ev);

    this.ws.onerror = (ev: Event): any => {
      this.toasts.push(ev.toString(), 'WS Request Error', 'error', 10000);
      this.emit('error', ev);
    };

    this.ws.onopen = (ev: Event): any => {
      this.emit('open', ev);
      this.sendMessage(WSCommand.INIT, {
        user_id: getCookieValue('userid'),
        token: getCookieValue('token'),
      });
    };

    this.ws.onclose = (ev: Event): any => {
      this.emit('close', ev);
      this._isInitialized = false;
      this.ws = new WebSocket(this.rootURL);
    };

    this.on(WSEvent.HELLO, () => (this._isInitialized = true));
  }

  private onMessage(ev: MessageEvent) {
    const msg = JSON.parse(ev.data) as Message;

    if (msg.name == WSEvent.ERROR) {
      this.toasts.push(
        (msg.data as WSErrorEvent).message,
        'WS Request Error',
        'error',
        10000
      );
      this.emit('error', msg.data);
      return;
    }

    this.emit(msg.name, msg.data);
  }

  public sendMessageRaw(payload: string) {
    this.ws.send(payload);
  }

  public sendMessage(name: WSCommand, data: any = null) {
    let payload = { name };
    if (data != undefined && data != null) {
      payload['data'] = data;
    }
    this.sendMessageRaw(JSON.stringify(payload));
  }

  public get isInitialized() {
    return this._isInitialized;
  }
}
