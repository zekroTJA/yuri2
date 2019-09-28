/** @format */

import { Component, forwardRef } from '@angular/core';
import { NG_VALUE_ACCESSOR, ControlValueAccessor } from '@angular/forms';

export const CUSTOM_INPUT_CONTROL_VALUE_ACCESSOR: any = {
  provide: NG_VALUE_ACCESSOR,
  // tslint:disable-next-line: no-use-before-declare
  useExisting: forwardRef(() => SliderComponent),
  multi: true,
};

@Component({
  selector: 'app-slider',
  templateUrl: './slider.component.html',
  styleUrls: ['./slider.component.sass'],
})
export class SliderComponent implements ControlValueAccessor {
  private _value: number;

  private onTouchedCallback: () => void = () => {};
  private onChangeCallback: (_: any) => void = () => {};

  public onChange(event: any) {
    if (!event || !event.target) {
      return;
    }

    this.value = event.target.value;
    console.log(this.value);
  }

  public get value(): number {
    return this._value;
  }

  public set value(v: number) {
    if (v !== this._value) {
      this._value = v;
      this.onChangeCallback(v);
    }
  }

  public onBlur() {
    this.onTouchedCallback();
  }

  public writeValue(v: number): void {
    if (v !== this._value) {
      this._value = v;
    }
  }

  public registerOnChange(fn: any): void {
    this.onChangeCallback = fn;
  }

  public registerOnTouched(fn: any): void {
    this.onTouchedCallback = fn;
  }
}
