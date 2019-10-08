/** @format */

import { Component, forwardRef, Output, EventEmitter } from '@angular/core';
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
  providers: [CUSTOM_INPUT_CONTROL_VALUE_ACCESSOR],
})
export class SliderComponent implements ControlValueAccessor {
  private _value: number;

  public ttVisible;
  public ttTranslation;

  @Output() public update: EventEmitter<any> = new EventEmitter();

  private onTouchedCallback: () => void = () => {};
  private onChangeCallback: (_: any) => void = () => {};

  constructor() {}

  public onChange(event: any) {
    if (!event || !event.target) {
      return;
    }

    this.value = event.target.valueAsNumber;
    this.update.emit(this.value);
  }

  public onInput(event: any) {
    if (!event || !event.target) {
      return;
    }

    this.value = event.target.valueAsNumber;

    this.setTranslation();
  }

  public onMouseIn() {
    this.ttVisible = true;
    this.setTranslation();
  }

  public onMouseOut() {
    this.ttVisible = false;
  }

  private setTranslation() {
    const maxSteps = 200;
    const transMin = -14;
    const transMax = 148;
    const trans = Math.floor(
      (this.value / maxSteps) * (transMax - transMin) + transMin
    );
    console.log(this.value);
    this.ttTranslation = `translate(${trans}px, -32px)`;
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
