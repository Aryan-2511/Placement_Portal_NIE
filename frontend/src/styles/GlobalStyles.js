import { createGlobalStyle } from 'styled-components';

const GlobalStyles = createGlobalStyle`

:root{
  &, &.light-mode{
    --color-grey-0 : #ffffff;
    --color-grey-50 : #f3f3f3;
    --color-grey-100 : #727272;
    --color-grey-600 : #3b3b3b;

    --color-blue-0 : #E5E5FF;
    --color-blue-100 : #B3DFF8;
    --color-blue-700 : #242466;

    --color-brown-100: #F8D1B3;
    --color-brown-700: #4B290F;

    --color-green-100: #BFFBD9;
    --color-green-700: #0B4324;

    --backdrop-color: rgba(255, 255, 255, 0.1);

    --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.04);
    --shadow-md: 0px 0.6rem 2.4rem rgba(0, 0, 0, 0.06);
    --shadow-lg: 0 2.4rem 3.2rem rgba(0, 0, 0, 0.12);

  }

  &.dark-mode{
    --color-grey-0 : #0B0916;
    --color-grey-50 : #14121F;
    --color-grey-100 : #AAAAAA;
    --color-grey-600 : #C4C4C4;

    --color-white : #fff;

    --color-blue-0 : #041524;
    --color-blue-100 : #242466;
    --color-blue-700 : #B3DFF8;

    --color-brown-100: #4B290F;
    --color-brown-700: #F8D1B3;

    --color-green-100: #0B4324;
    --color-green-700: #BFFBD9;
  }
}


*,
*::before,
*::after {
  box-sizing: border-box;
  padding: 0;
  margin: 0;

  /* Creating animations for dark mode */
  transition: background-color 0.3s, border 0.3s;
}
html{
  font-size: 62.5%;
}
body {
  font-family: "Oxygen", sans-serif;
  /* color: var(--color-grey-700); */

  transition: color 0.3s, background-color 0.3s;
  min-height: 100vh;
  line-height: 1.5;
  font-size: 1.6rem;
  background-color: var(--color-grey-50);
}

button {
  cursor: pointer;
}

a {
  color: inherit;
  text-decoration: none;
}
p{
  color: var(--color-grey-600)
}
ul {
  list-style: none;
}

p,
h1,
h2,
h3,
h4,
h5,
h6 {
  overflow-wrap: break-word;
  hyphens: auto;
}`;

export default GlobalStyles;
