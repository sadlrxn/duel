import React, { FC } from 'react';
import Select, { Options, StylesConfig, Props } from 'react-select';

interface CustomSelectProps extends Props {
  background: string;
  hoverBackground: string;
  color: string;
  fontSize?: string;
  fontWeight?: number;
  width?: string | number;
  minHeight?: string | number;
  paddingLeft?: string;
  borderRadius?: string | number;
  options: Options<any>;
}
const CustomSelect: FC<CustomSelectProps> = ({
  options,
  background,
  hoverBackground,
  color,
  fontSize = '16px',
  fontWeight = 600,
  width,
  minHeight,
  paddingLeft,
  borderRadius,
  ...props
}) => {
  const customStyles: StylesConfig = {
    container: provided => ({
      ...provided,
      width: width ?? provided.width
    }),
    control: provided => ({
      ...provided,
      background,
      border: 0,
      minHeight: minHeight ?? provided.minHeight,
      paddingLeft: paddingLeft ?? provided.paddingLeft,
      borderRadius: borderRadius ?? provided.borderRadius,
      boxShadow: 'none',
      cursor: 'pointer',
      '&:hover': {
        background: hoverBackground
      }
    }),
    option: provided => ({
      ...provided,
      background,
      fontFamily: 'Inter',
      fontWeight,
      fontSize,
      color,
      cursor: 'pointer',
      '&:hover': {
        background: hoverBackground
      }
    }),
    input: base => ({
      ...base
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
      fontWeight,
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
