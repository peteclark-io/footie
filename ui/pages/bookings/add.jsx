import Link from 'next/link'
import Router from 'next/router'
import { withRouter } from 'next/router'
import React from 'react'
import axios from 'axios'
import {Typography} from '@material-ui/core'
import {TextField, Button} from '@material-ui/core';

import Breadcrumbs from '../../components/breadcrumbs'
import Progress from '../../components/progress'

class AddBooking extends React.Component {

  constructor(props) {
    super(props);
    this.state = {loading: true};

    this.handleChange = this.handleChange.bind(this);
    this.submitForm = this.submitForm.bind(this);
  }

  componentDidMount(){
    axios.get('http://localhost:8000/groups/' + this.props.router.query.forGroup)
      .then((resp) => {
        this.setState({group: resp.data, loading: false})
      });
    this.setState({loading: true})
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
    axios.post('http://localhost:8000/bookings', {displayName: this.state.displayName, email: this.state.email})
      .then((resp) => {
        var timeout = setTimeout(() => Router.push(`/bookings/${resp.data.id}`), 1000)
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
        <Typography variant="h6" gutterBottom={true}>Booking created successfully!</Typography>
      )
    }

    const crumbs = [{title: 'Group'}, {title: this.state.group.name, href: `/groups?id=${this.state.group.id}`}]

    return (
      <section>
        <Breadcrumbs crumbs={crumbs}></Breadcrumbs>
        <Typography variant="h5">Add New Booking</Typography>

        <section>
          <form className={'container form'} onSubmit={this.submitForm}>
            <div>
              <TextField id="start"
                required
                type="date"
                name="start"
                margin="normal"
                helperText="When does the booking begin?"
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
              Add Booking
            </Button>
          </form>

        </section>

      </section>
    )
  }
}

export default withRouter(AddBooking);
