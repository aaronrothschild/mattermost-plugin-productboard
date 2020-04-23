import {getConfig} from 'mattermost-redux/selectors/entities/general';
import {Client4} from 'mattermost-redux/client';
import {id as pluginId} from './manifest';


export const getPluginServerRoute = (state) => {
    const config = getConfig(state);

    let basePath = '';
    if (config && config.SiteURL) {
        basePath = new URL(config.SiteURL).pathname;

        if (basePath && basePath[basePath.length - 1] === '/') {
            basePath = basePath.substr(0, basePath.length - 1);
        }
    }

    return basePath + '/plugins/' + pluginId;
};


export const create = (postID) => async (dispatch, getState) => {
    {
        await fetch(getPluginServerRoute(getState()) + '/create', Client4.getOptions({
            method: 'post',
            body: JSON.stringify({post_id: postID}),
        }));
    }
};