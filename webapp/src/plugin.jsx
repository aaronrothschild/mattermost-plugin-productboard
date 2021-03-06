import React from 'react';

import {FormattedMessage} from 'react-intl';

// import en from 'i18n/en.json';

// import es from 'i18n/es.json';

import {id as pluginId} from './manifest';
import { create } from './actions';

// import Root from './components/root';
// import BottomTeamSidebar from './components/bottom_team_sidebar';
// import LeftSidebarHeader from './components/left_sidebar_header';
// import LinkTooltip from './components/link_tooltip';
// import UserAttributes from './components/user_attributes';
// import UserActions from './components/user_actions';
// import RHSView from './components/right_hand_sidebar';
// import SecretMessageSetting from './components/admin_settings/secret_message_setting';
// import CustomSetting from './components/admin_settings/custom_setting';

// import PostType from './components/post_type';
// import EphemeralPostType from './components/ephemeral_post_type';
// import {
//     MainMenuMobileIcon,
//     ChannelHeaderButtonIcon,
//     FileUploadMethodIcon,
// } from './components/icons';
// import {
//     mainMenuAction,
//     fileUploadMethodAction,
//     postDropdownMenuAction,
//     postDropdownSubMenuAction,
//     channelHeaderMenuAction,
//     websocketStatusChange,
//     getStatus,
// } from './actions';
// import reducer from './reducer';

// function getTranslations(locale) {
//     switch (locale) {
//     case 'en':
//         return en;
//     case 'es':
//         return es;
//     }
//     return {};
// }




export default class Plugin {
    initialize(registry, store) {
        registry.registerPostDropdownMenuAction(
            <FormattedMessage
                id='com.mattermost.productboard'
                defaultMessage='Send to ProductBoard'
            />,
            (postID) => store.dispatch(create(postID)),
        );


        // // Immediately fetch the current plugin status.
        // store.dispatch(getStatus());

        // // Fetch the current status whenever we recover an internet connection.
        // registry.registerReconnectHandler(() => {
        //     store.dispatch(getStatus());
        // });

        // registry.registerTranslations(getTranslations);
    }

    uninitialize() {
        //eslint-disable-next-line no-console
        console.log(pluginId + '::uninitialize()');
    }
}