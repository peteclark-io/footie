import Link from 'next/link'
import Router from 'next/router'
import { withRouter } from 'next/router'
import React from 'react'
import axios from 'axios'
import {Typography} from '@material-ui/core'
import {TextField, Button} from '@material-ui/core';

import Progress from '../../components/progress'

class AddPlayer extends React.Component {

  constructor(props) {
    super(props);
    this.state = {};

    this.handleChange = this.handleChange.bind(this);
    this.submitForm = this.submitForm.bind(this);
  }

  componentWillUnmount(){
    if(this.state.timeout){
      clearTimeout(this.state.timeout)
    }
  }

  handleChange(e){
    this.setState({[e.target.id]: e.target.value})
  }

  submitForm(e){
    e.preventDefault();
    axios.post('http://localhost:8000/players', {displayName: this.state.displayName, email: this.state.email})
      .then((resp) => {
        var timeout = setTimeout(() => Router.push(`/players/${resp.data.id}`), 1000)
        this.setState({submitted: {ok: true}, loading: false, timeout: timeout})
      });
    this.setState({loading: true})
  }

  render(){
    if(this.state.loading){
      return (
        <Progress></Progress>
      )
    }

    if(this.state.submitted){
      return (
        <Typography variant="h6" gutterBottom={true}>Player created successfully!</Typography>
      )
    }

    return (
      <section>
        <Typography variant="h4" gutterBottom={true}>Add New Player</Typography>

        <section>
          <form className={'container form'} onSubmit={this.submitForm}>
            <div>
              <TextField id="displayName"
                required
                label="Name"
                type="text"
                name="name"
                margin="normal"
                helperText="Display name"
                onChange={this.handleChange} />
            </div>

            <div className='sm-margin-bottom'>
              <TextField id="email"
                required
                label="Email"
                type="email"
                name="email"
                autoComplete="email"
                margin="normal"
                helperText="A valid email address"
                onChange={this.handleChange} />
            </div>

            <Button variant="contained" color="secondary" type="submit">
              Add Player
            </Button>
          </form>

        </section>

      </section>
    )
  }
}

export default withRouter(AddPlayer);
