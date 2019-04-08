import React, {Component} from 'react';
import logo from './DoGo.jpg';

class HomeScreen extends Component {
    render(){
        return(
            <div style={{minHeight: '100%', alignContent: 'center'}}>
                <div style={{width: '100%', alignContent:'center'}}>
                    <img src={logo} style={{ display: 'block',margin: 'auto', marginTop: 40,}}></img>
                </div>
                
            </div>
        );
    }
}

export default HomeScreen;