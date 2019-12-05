import Router from 'next/router'
import { withRouter } from 'next/router'
import React from 'react'
import axios from 'axios'
import {Typography} from '@material-ui/core'

import Progress from '../../components/progress'

class Player extends React.Component {

  constructor(props) {
    super(props);
    this.state = {loading: true};
  }

  componentDidMount(){
    axios.get('http://localhost:8000/players/' + this.props.router.query.id)
      .then((resp) => {
        this.setState({loading: false, ...resp.data})
      });
  }

  renderDataPair(key, value){
    if(!value){
      return (
        <React.Fragment></React.Fragment>
      )
    }

    return (
      <React.Fragment>
        <Typography variant="h6" gutterBottom={true}>{key}</Typography>
        <Typography>{value}</Typography>
      </React.Fragment>
    )
  }

  render(){
    if(this.state.loading){
      return (
        <Progress></Progress>
      )
    }

    const {displayName, email} = this.state;

    return (
      <section>
        <Typography className="sm-margin-bottom" variant="h4" gutterBottom={true}>{displayName}</Typography>
        {this.renderDataPair('Email', email)}
      </section>
    )
  }
}

export default withRouter(Player);
