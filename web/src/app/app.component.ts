/** @format */

import { Component } from '@angular/core';
import { WSService } from './api/ws/ws.service';
import { ToastService } from './components/toast/toast.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.sass'],
})
export class AppComponent {
  title = 'web';
  constructor(private ws: WSService, public toasts: ToastService) {}
}
