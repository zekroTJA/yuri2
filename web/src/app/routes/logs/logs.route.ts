/** @format */

import { Component } from '@angular/core';
import { RestService } from 'src/app/api/rest/rest.service';
import { LogEntry } from 'src/app/api/rest/rest.models';
import { SharedDataService } from 'src/app/services/shareddata.service';
import { ToastService } from 'src/app/components/toast/toast.service';
import { WSService } from 'src/app/api/ws/ws.service';
import dateFormat from 'dateformat';
import { SourceType } from 'src/app/api/ws/ws.static';

@Component({
  selector: 'app-logs-route',
  templateUrl: './logs.route.html',
  styleUrls: ['./logs.route.sass'],
})
export class LogsRouteComponent {
  public log: LogEntry[];
  public dateFormat = dateFormat;

  constructor(
    private rest: RestService,
    private sharedData: SharedDataService,
    private ws: WSService,
    private toasts: ToastService
  ) {
    this.fetchLogs();
  }

  public fetchLogs() {
    if (!this.ws.isInitialized) {
      setTimeout(this.fetchLogs.bind(this), 100);
      return;
    }

    if (!this.sharedData.currentGuildID) {
      this.toasts.push(
        "You need to be in a Voice Channel with Yuri that the Guild's log can be fetched!",
        'Error',
        'error',
        10000,
        true
      );
      return;
    }

    this.rest
      .getPlayLog(this.sharedData.currentGuildID, 0, 100)
      .subscribe((log) => (this.log = log));
  }

  public getSourceClassName(source: string) {
    return `source-${source.toLowerCase()}`;
  }
}
