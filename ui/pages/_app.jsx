import React from "react";

import App, { Container } from 'next/app'
import Head from 'next/head'

import JssProvider from 'react-jss/lib/JssProvider'
import CssBaseline from '@material-ui/core/CssBaseline'
import { MuiThemeProvider } from '@material-ui/core/styles'
import getPageContext from '../core/material-ui'

import '../styles/html.scss'
import Header from '../components/header'

class CustomApp extends App {

  constructor() {
    super();
    this.pageContext = getPageContext();
  }

  componentDidMount() {
    // Remove the server-side injected CSS.
    const jssStyles = document.querySelector('#jss-server-side');
    if (jssStyles && jssStyles.parentNode) {
      jssStyles.parentNode.removeChild(jssStyles);
    }
  }

  static async getInitialProps({ Component, ctx }) {
    let pageProps = {};

    if (Component.getInitialProps) {
      pageProps = await Component.getInitialProps(ctx);
    }

    return { pageProps };
  }

  render() {
    const { Component, pageProps } = this.props;

    return (
      <Container>
        <JssProvider registry={this.pageContext.sheetsRegistry} generateClassName={this.pageContext.generateClassName}>
          <MuiThemeProvider theme={this.pageContext.theme} sheetsManager={this.pageContext.sheetsManager}>
            <CssBaseline></CssBaseline>

            <Header></Header>
            <section className="entry">
              <Component pageContext={this.pageContext} {...pageProps} />
            </section>
          </MuiThemeProvider>
        </JssProvider>
      </Container>
    );
  }
}

export default CustomApp;
