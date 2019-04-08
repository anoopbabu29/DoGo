import React,  {cockpit,Component} from 'react';
import AppBar from '@material-ui/core/AppBar';
import Toolbar from '@material-ui/core/Toolbar';
import Typography from '@material-ui/core/Typography';
import Paper from '@material-ui/core/Paper';
import InputBase from '@material-ui/core/InputBase';
import Divider from '@material-ui/core/Divider';
import IconButton from '@material-ui/core/IconButton';
import MenuIcon from '@material-ui/icons/Menu';
import SearchIcon from '@material-ui/icons/Search';
import Chip from '@material-ui/core/Chip';
import pic from './DoGo Logo@300x.png';
import CatBar from './subsections/CatBar';
import $ from 'jquery';

function handleClick() {
    alert('You clicked the Chip.'); // eslint-disable-line no-alert
}

const style = {
    chipLabel:{
        backgroundColor: 'white', 
        marginBottom: 10,
        marginRight: 5,
    }
}

function getQueries(){
    $.ajax({url: "https://localhost:3000/api/getQuery", method: "POST", data: $('#inp').val()}).done(function(resp){console.log(resp)});
    
}

class NavbarGo extends Component {
    render(){
        return(
            <div>
                {/* Navbar */}
                <AppBar  style={{backgroundColor: '#73CEDD', height: 75}}>
                    <Toolbar style={{marginTop: 5,}}>
                        <img src={pic}  height="50" width="120"/>

                        {/* Search Bar */}
                        <div style={{alignContent: 'center', marginLeft: 50}}>
                            <Paper style={{padding: '2px 4px',display: 'flex',alignItems: 'center',width: 600,marginLeft: 'auto', marginRight: 'auto', borderRadius: 50}}  elevation={0}>   
                                <InputBase id="inp" placeholder="Search to Donate" style={{width: 600, marginLeft: 15}} />
                                <Divider styles={{width: 1,height: 28,margin: 4,}} />
                                <IconButton color="primary"  >
                                    <SearchIcon />
                                </IconButton>
                            </Paper>
                        </div>
                    </Toolbar>
                    
                </AppBar>
                <Paper  style={{backgroundColor: 'lightgray', opacity:5, marginTop: 75, height: 50, borderRadius: 0}}>
                    <Toolbar id="chipContainer" style={{}}>
                        <Chip
                            label="Arts, Culture and Humanities"
                            onClick={handleClick}
                            variant="outlined"
                            style={style.chipLabel}
                        />
                        <Chip
                            label="Education and Research"
                            onClick={handleClick}
                            variant="outlined"
                            style={style.chipLabel}
                        />
                        <Chip
                            label="Environment and Animals"
                            onClick={handleClick}
                            variant="outlined"
                            style={style.chipLabel}
                        />
                        <Chip
                            label="Health"
                            onClick={handleClick}
                            variant="outlined"
                            style={style.chipLabel}
                        />
                        <Chip
                            label="Human Services"
                            onClick={handleClick}
                            variant="outlined"
                            style={style.chipLabel}
                        />
                        <Chip
                            label="International"
                            onClick={handleClick}
                            variant="outlined"
                            style={style.chipLabel}
                        />
                        <Chip
                            label="Public, Societal Benefit"
                            onClick={handleClick}
                            variant="outlined"
                            style={style.chipLabel}
                        />
                        <Chip
                            label="Religion"
                            onClick={handleClick}
                            variant="outlined"
                            style={style.chipLabel}
                        />
                        <Chip
                            label="Other"
                            onClick={handleClick}
                            variant="outlined"
                            style={style.chipLabel}
                        />
                    </Toolbar>
                </Paper>
            </div>
        );
    }
}

export default NavbarGo;