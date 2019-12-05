import Link from 'next/link'
import { withRouter } from 'next/router'
import React from 'react'
import axios from 'axios'
import {Typography} from '@material-ui/core'
import {Table, TableBody, TableCell, TableHead, TableRow} from '@material-ui/core'

class Match extends React.Component {

  state = {}

  componentDidMount(){
    axios.get('http://localhost:8000/groups/' + this.props.router.query.id)
      .then((resp) => {
        this.setState(resp.data)
      });
  }

  render(){
    const {name} = this.state;

    if(!name){
      return <React.Fragment></React.Fragment>
    }

    return (
      <section>
        <Typography className="sm-margin-bottom" variant="h4" gutterBottom={true}>{name}</Typography>
      </section>
    )
  }
}

export default withRouter(Match);
