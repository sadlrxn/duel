import React from 'react';
import { StylesConfig, Options, Props, createFilter } from 'react-select';
import Creatable from 'react-select/creatable';

export interface CreatableSelectProps extends Props {
  background: string;
  hoverBackground: string;
  color: string;
  fontSize?: string;
  width?: number | string;
  maxWidth?: number | string;
  options: Options<{
    label: string;
    value: string;
  }>;
  onCreateOption?: any;
}

export default function CreatableSelect({
  background = '',
  hoverBackground = '',
  color = '',
  fontSize = '16px',
  width = 'auto',
  maxWidth = 'auto',
  options,
  onCreateOption,
  ...props
}: CreatableSelectProps) {
  const customStyles: StylesConfig = {
    control: provided => ({
      ...provided,
      background,
      border: 0,
      width,
      maxWidth,
      boxShadow: 'none',
      cursor: 'pointer',
      '&:hover': {
        background: hoverBackground
      }
    }),
    option: (base, { isFocused }) => ({
      ...base,
      fontFamily: 'Inter',
      fontWeight: '600',
      fontSize,
      color,
      cursor: 'pointer',
      background: isFocused ? hoverBackground : background
      // '&:hover': {
      //   background: hoverBackground
      // }
    }),
    input: base => ({
      ...base,
      fontSize,
      color,
      fontWeight: '600',
      fontFamily: 'Inter'
    }),

    singleValue: provided => ({
      ...provided,
      color
    }),
    indicatorSeparator: () => ({ display: 'none' }),
    dropdownIndicator: (provided, state) => ({
      ...provided,
      color,
      transition: '0.5s',
      transform: state.selectProps.menuIsOpen
        ? 'rotate(180deg)'
        : 'rotate(0deg)',
      '&:hover': {
        color
      }
    }),
    menu: provided => ({
      ...provided,
      background: 'transparent'
    }),
    menuList: provided => ({
      ...provided,
      background,
      borderRadius: '7px'
    }),
    valueContainer: base => ({
      ...base,
      fontFamily: 'Inter',
      fontWeight: '600',
      fontSize,
      color,
      minWidth: '100px'
    })
  };

  return (
    <Creatable
      options={options}
      styles={customStyles}
      defaultValue={options[0]}
      onCreateOption={onCreateOption}
      filterOption={createFilter({ ignoreCase: false })}
      //@ts-ignore
      isValidNewOption={(inputValue, selectValue, selectOptions) => {
        const exactValueExists = selectOptions.find(
          //@ts-ignore
          el => el.value === inputValue
        );
        const valueIsNotEmpty = inputValue.trim().length;
        return !exactValueExists && valueIsNotEmpty;
      }}
      {...props}
    />
  );
}
