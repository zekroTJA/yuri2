/** @format */

import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import { WSService } from 'src/app/api/ws/ws.service';
import { WSCommand, SourceType } from 'src/app/api/ws/ws.static';
import { DomSanitizer, SafeResourceUrl } from '@angular/platform-browser';
import { RestService } from 'src/app/api/rest/rest.service';
import { ToastService } from '../toast/toast.service';
import { YouTubeEmbed } from 'src/app/api/rest/rest.models';

@Component({
  selector: 'app-ytimport',
  templateUrl: './ytimport.component.html',
  styleUrls: ['./ytimport.component.sass'],
})
export class YTInputComponent implements OnInit, OnDestroy {
  public ident: string;
  public resourceURL: SafeResourceUrl;
  public videoInfo: YouTubeEmbed;

  public rendered: boolean;
  public visible: boolean;

  constructor(
    private ws: WSService,
    private rest: RestService,
    private toasts: ToastService,
    private snaitazier: DomSanitizer
  ) {}

  public onCancel() {
    this.close();
  }

  public onPlay() {
    this.ws.sendMessage(WSCommand.PLAY, {
      ident: this.ident,
      source: SourceType.YOUTUBE,
    });

    this.close();
  }

  private close() {
    this.visible = false;
    setTimeout(() => (this.rendered = false), 500);
  }

  private show() {
    this.rendered = true;
    setTimeout(() => (this.visible = true), 10);
  }

  private onKeyDown(event: any) {
    if (event.keyCode == 27) {
      this.close();
      event.preventDefault();
    }
  }

  private onPaste(event: any) {
    event.stopPropagation();
    event.preventDefault();

    let pastedData: string = event.clipboardData.getData('Text');
    let i: number;
    let cut: number;

    if (!pastedData) return;

    i = pastedData.indexOf('&');
    if (i > -1) pastedData = pastedData.substring(0, i);

    i = pastedData.indexOf('?');
    if (i > -1 && pastedData.substr(i - 5, 5) !== 'watch')
      pastedData = pastedData.substring(0, i);

    i = pastedData.indexOf('youtube.com/watch?v=');
    if (i > -1) cut = i + 'youtube.com/watch?v='.length;

    i = pastedData.indexOf('youtu.be/');
    if (i > -1) cut = i + 'youtu.be/'.length;

    pastedData = pastedData.substring(cut);

    this.rest
      .getYouTubeEmbed('https://youtu.be/' + pastedData)
      .toPromise()
      .then((data: YouTubeEmbed) => {
        if (data.error) {
          this.toasts.push(
            'The pasted URL does not refer to a valid or accessable YouTube video!',
            'Invalid YouTube Video',
            'error',
            10000,
            true
          );
          return;
        }

        this.ident = pastedData;
        this.videoInfo = data;

        // this.resourceURL = this.snaitazier.bypassSecurityTrustResourceUrl(
        //   'https://www.youtube-nocookie.com/embed/' + pastedData
        // );

        this.show();
      })
      .catch((e) => {
        this.toasts.push(
          `Error on validating YouTube video: ${e}`,
          'Validation Error',
          'error',
          10000,
          true
        );
      });
  }

  public ngOnInit() {
    window.addEventListener('keydown', this.onKeyDown.bind(this));
    window.addEventListener('paste', this.onPaste.bind(this));
  }

  public ngOnDestroy() {
    window.removeEventListener('keydown', this.onKeyDown.bind(this));
    window.removeEventListener('paste', this.onPaste.bind(this));
  }
}
