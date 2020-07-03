import React, { Component } from 'react';
import {
    AppBar,
    Divider,
    Drawer,
    IconButton,
    List,
    ListItem,
    ListItemIcon,
    ListItemText,
    MuiThemeProvider,
    Toolbar,
    Typography
} from '@material-ui/core';
import MenuIcon from '@material-ui/icons/Menu';
import SettingsIcon from '@material-ui/icons/Settings';
import ChevronLeftIcon from '@material-ui/icons/ChevronLeft';
import VerticalSplitIcon from '@material-ui/icons/VerticalSplit';
import PhotoIcon from '@material-ui/icons/Photo';
import PhotoAlbumIcon from '@material-ui/icons/PhotoAlbum';
import clsx from 'clsx';
import history from '../../Routing/History';
import './TitleBar.css';
import theme from '../../Theme/theme';

interface IProps {
  title: string;
}

interface IState {
  menuIsOpen: boolean;
}

export class TitleBar extends Component<IProps, IState> {

  constructor(props: IProps) {
    super(props);
    this.state = { menuIsOpen: false };
  }

  public siteLinks = [{
    name: 'Dashboard',
    icon: <VerticalSplitIcon />,
    link: '/',
  },
  {
    name: 'Photos',
    icon: <PhotoIcon />,
    link: '/photos'
  },
  {
    name: 'Albums',
    icon: <PhotoAlbumIcon />,
    link: '/albums',
  }];

  public metaLinks = [{
    name: 'Settings',
    icon: <SettingsIcon />,
    link: '/settings',
  }];

  public render() {
    return(
            <div className="titleBar__mainBar">
                <MuiThemeProvider theme={theme}>
                    <AppBar
                        position="fixed"
                        className="titleBar__appBar"
                    >
                        <Toolbar>
                            <IconButton
                                color="inherit"
                                aria-label="open drawer"
                                onClick={() => this.setState({ menuIsOpen: true })}
                                edge="start"
                                className={clsx('titleBar__menuButton', this.state.menuIsOpen && 'titleBar__iconButtonHide')}
                            >
                                <MenuIcon />
                            </IconButton>
                            <Typography variant="h6">
                                {this.props.title}
                            </Typography>
                            {this.props.children}
                        </Toolbar>
                    </AppBar>
                    <Drawer
                        className="titleBar__drawer"
                        variant="persistent"
                        anchor="left"
                        open={this.state.menuIsOpen}
                        classes={{ paper: 'titleBar__drawerPaper' }}
                    >
                        <div className="titleBar__drawerHeader">
                            <div className="titleBar__drawerHeader__title">Menu</div>
                            <IconButton onClick={() => this.setState({ menuIsOpen: false })}>
                                <ChevronLeftIcon />
                            </IconButton>
                        </div>
                        <Divider />
                        <List>
                            {/* tslint:disable-next-line:jsx-no-multiline-js */}
                            {this.siteLinks.map((link, index) => (
                                <ListItem
                                    button={true}
                                    key={link.name}
                                    onClick={(() => history.push(link.link))}
                                >
                                    <ListItemIcon>{link.icon}</ListItemIcon>
                                    <ListItemText primary={link.name} />
                                </ListItem>
                            ))}
                        </List>
                        <Divider />
                        <List>
                            {/* tslint:disable-next-line:jsx-no-multiline-js */}
                            {this.metaLinks.map((link, index) => (
                                <ListItem
                                    button={true}
                                    key={link.name}
                                    onClick={(() => history.push(link.link))}
                                >
                                    <ListItemIcon>{link.icon}</ListItemIcon>
                                    <ListItemText primary={link.name} />
                                </ListItem>
                            ))}
                        </List>
                    </Drawer>
                </MuiThemeProvider>
            </div>
    );
  }
}
