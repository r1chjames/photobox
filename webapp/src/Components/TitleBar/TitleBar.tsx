import React from 'react';
import {
  AppBar,
  Divider,
  Drawer,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  MuiThemeProvider,
  Toolbar, Typography
} from '@material-ui/core';
import SettingsIcon from '@material-ui/icons/Settings';
import VerticalSplitIcon from '@material-ui/icons/VerticalSplit';
import PhotoIcon from '@material-ui/icons/Photo';
import PhotoAlbumIcon from '@material-ui/icons/PhotoAlbum';
import history from '../../Routing/History';
import './TitleBar.css';
import theme from '../../Theme/theme';

export const TitleBar: React.FunctionComponent = (props) => {

  const siteLinks = [{
    name: 'Dashboard',
    icon: <VerticalSplitIcon/>,
    link: '/',
  },
    {
      name: 'Photos',
      icon: <PhotoIcon/>,
      link: '/photos'
    },
    {
      name: 'Albums',
      icon: <PhotoAlbumIcon/>,
      link: '/albums',
    }];

  const metaLinks = [{
    name: 'Settings',
    icon: <SettingsIcon/>,
    link: '/settings',
  }];

  return (
        <div className="titleBar__mainBar">
            <MuiThemeProvider theme={theme}>
                <AppBar
                    position="fixed"
                    className="titleBar__appBar"
                    elevation={0}
                >
                    <Toolbar>
                      <Typography variant="h6">
                        Photobox
                      </Typography>
                        {props.children}
                    </Toolbar>
                    <Divider/>
                </AppBar>
                <Drawer
                    className="titleBar__drawer"
                    variant="permanent"
                    anchor="left"
                    classes={{ paper: 'titleBar__drawerPaper' }}
                >
                    <Divider/>
                    <List>
                        {/* tslint:disable-next-line:jsx-no-multiline-js */}
                        {siteLinks.map((link, index) => (
                            <ListItem
                                button={true}
                                key={link.name}
                                onClick={(() => history.push(link.link))}
                            >
                                <ListItemIcon>{link.icon}</ListItemIcon>
                                <ListItemText primary={link.name}/>
                            </ListItem>
                        ))}
                    </List>
                    <List>
                        {/* tslint:disable-next-line:jsx-no-multiline-js */}
                        {metaLinks.map((link, index) => (
                            <ListItem
                                button={true}
                                key={link.name}
                                onClick={(() => history.push(link.link))}
                            >
                                <ListItemIcon>{link.icon}</ListItemIcon>
                                <ListItemText primary={link.name}/>
                            </ListItem>
                        ))}
                    </List>
                </Drawer>
            </MuiThemeProvider>
        </div>
  );
};
