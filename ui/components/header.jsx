import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';

const Header = () => (
  <div className="grow">
    <AppBar position="static">
      <Toolbar>
        <IconButton className="" color="inherit" aria-label="Menu">
          <MenuIcon />
        </IconButton>
        <Typography variant="h6" color="inherit" className="grow">
          EL Football
        </Typography>
      </Toolbar>
    </AppBar>
  </div>
);

export default Header
