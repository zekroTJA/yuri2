/** @format */

import {
  Component,
  Input,
  TemplateRef,
  ElementRef,
  ViewChild,
} from '@angular/core';

export interface ContextMenuItem {
  el: string | TemplateRef<any>;
  action: () => void;
}

@Component({
  selector: 'app-contextmenu',
  templateUrl: './contextmenu.component.html',
  styleUrls: ['./contextmenu.component.sass'],
})
export class ContextMenuComponent {
  @ViewChild('contextmenu', { static: false }) element: ElementRef;

  @Input() x: number;
  @Input() y: number;
  @Input() visible: boolean;

  @Input() items: ContextMenuItem[] = [];

  public isTemplate(element) {
    return element.el instanceof TemplateRef;
  }
}
