import Link from 'next/link'
import { withRouter } from 'next/router'
import React from 'react'
import axios from 'axios'
import {Typography} from '@material-ui/core'
import {Table, TableBody, TableCell, TableHead, TableRow} from '@material-ui/core'

class Booking extends React.Component {

  state = {}

  componentDidMount(){
    axios.get('http://localhost:8000/bookings/' + this.props.router.query.id)
      .then((resp) => {
        this.setState(resp.data)
      });
  }

  render(){
    const {group, start} = this.state;

    if(!name){
      return <React.Fragment></React.Fragment>
    }

    return (
      <section>
        <Typography className="sm-margin-bottom" variant="h4" gutterBottom={true}>{name}</Typography>

        <Typography variant="h5">Player Roster</Typography>

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

export default withRouter(Booking);
