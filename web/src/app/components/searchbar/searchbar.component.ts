/** @format */

import {
  Component,
  EventEmitter,
  ViewChild,
  ElementRef,
  AfterViewInit,
  Output,
} from '@angular/core';

@Component({
  selector: 'app-searchbar',
  templateUrl: './searchbar.component.html',
  styleUrls: ['./searchbar.component.sass'],
})
export class SearchBarComponent implements AfterViewInit {
  @ViewChild('inptbar', { static: true }) inptbar: ElementRef;

  @Output() public search: EventEmitter<any> = new EventEmitter();
  @Output() public close: EventEmitter<any> = new EventEmitter();

  public ngAfterViewInit() {
    this.inptbar.nativeElement.focus();
  }

  public onCloseClick() {
    this.close.emit();
  }

  public onInput(ev: any) {
    const target = ev.target;
    if (!target) return;

    this.search.emit(target.value);
  }
}
