import React, { FC } from 'react';
import Select, { Options, StylesConfig, Props } from 'react-select';

interface CustomSelectProps extends Props {
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
}
const CustomSelect: FC<CustomSelectProps> = ({
  options,
  background,
  hoverBackground,
  width = 'auto',
  maxWidth = 'auto',
  color,
  fontSize = '16px',
  ...props
}) => {
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
    <Select
      defaultValue={options[0]}
      styles={customStyles}
      options={options}
      {...props}
      // components={{
      //   Menu: CustomMenu,
      // }}
    />
  );
};

export default CustomSelect;
