import Link from 'next/link'
import { withRouter } from 'next/router'
import React from 'react'
import axios from 'axios'
import {Typography} from '@material-ui/core'
import {Table, TableBody, TableCell, TableHead, TableRow} from '@material-ui/core'

import Breadcrumbs from '../../components/breadcrumbs'

class Group extends React.Component {

  state = {}

  componentDidMount(){
    axios.get('http://localhost:8000/groups/' + this.props.router.query.id)
      .then((resp) => {
        this.setState(resp.data)
      });
  }

  render(){
    const {name, players} = this.state;

    if(!name){
      return <React.Fragment></React.Fragment>
    }

    const crumbs = [{title: name}];

    return (
      <section>
        <Breadcrumbs crumbs={crumbs}></Breadcrumbs>

        <Typography variant="h6">Player Roster</Typography>

        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Name</TableCell>
              <TableCell>Email</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {players.map(player => (
              <TableRow key={player.id}>
                <TableCell component="th" scope="row">
                  {player.displayName}
                </TableCell>
                <TableCell>{player.email}</TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </section>
    )
  }
}

export default withRouter(Group);
