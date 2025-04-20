import { View, StyleSheet, LayoutChangeEvent } from 'react-native';
import { BottomTabBarProps } from '@react-navigation/bottom-tabs';
import TabBarButton from './TabBarButton';
import { Colors } from '@/constants/Colors';
import { useEffect, useState, useMemo } from 'react';
import { useAnimatedStyle, useSharedValue, withTiming } from 'react-native-reanimated';
import Animated from 'react-native-reanimated';
import React from 'react';

export default function TabBar({ state, descriptors, navigation }: BottomTabBarProps) {
    const [dimensions, setDimensions] = useState({ width: 0, height: 0 });

    // Map route names to their indices for direct access 
    const routeMap = useMemo(() => {
        const map: Record<string, number> = {};
        state.routes.forEach((route: {name: string}, index: number) => {
            map[route.name] = index;
        });
        return map;
    }, [state.routes]);

    const buttonWidth = dimensions.width / state.routes.length;
    const indicatorWidth = buttonWidth / 2;
    
    const getIndicatorPosition = (index: number) => {
        return (buttonWidth * index) + (buttonWidth / 2) - (indicatorWidth / 2);
    };

    const tabPositionX = useSharedValue(0);

    useEffect(() => {
        if (dimensions.width > 0) {
            const position = getIndicatorPosition(state.index);
            tabPositionX.value = withTiming(position, { duration: 200 });
        }
    }, [state.index, dimensions.width]);

    const onTabLayout = (event: LayoutChangeEvent) => {
        const { width, height } = event.nativeEvent.layout;
        if (width !== dimensions.width || height !== dimensions.height) {
            setDimensions({ width, height });
            
            if (width > 0) {
                const position = getIndicatorPosition(state.index);
                tabPositionX.value = position;
            }
        }
    };

    const animatedStyle = useAnimatedStyle(() => {
        return {
            transform: [
                {
                    translateX: tabPositionX.value
                }
            ]
        }
    });

    return (
        <View onLayout={onTabLayout} style={styles.tabBar}>
            <Animated.View
                style={[animatedStyle, {
                    position: "absolute",
                    backgroundColor: Colors.primary,
                    top: 0,
                    width: indicatorWidth,
                    height: 2,
                }]}
            />
            {state.routes.map((route: any, index: any) => {
                const { options } = descriptors[route.key];
                const label =
                    options.tabBarLabel !== undefined
                        ? options.tabBarLabel
                        : options.title !== undefined
                            ? options.title
                            : route.name;

                const isFocused = state.index === index;

                const onPress = () => {
                    const event = navigation.emit({
                        type: 'tabPress',
                        target: route.key,
                        canPreventDefault: true,
                    });

                    if (!isFocused && !event.defaultPrevented) {
                        navigation.navigate(route.name, route.params);
                    }
                };

                const onLongPress = () => {
                    navigation.emit({
                        type: 'tabLongPress',
                        target: route.key,
                    });
                };

                return (
                    <TabBarButton
                        key={route.name}
                        onPress={onPress}
                        onLongPress={onLongPress}
                        isFocused={isFocused}
                        label={label as string}
                        routeName={route.name}
                    />
                );
            })}
        </View>
    );
}

const styles = StyleSheet.create({
    tabBar: {
        flexDirection: 'row',
        paddingTop: 16,
        paddingBottom: 40,
        backgroundColor: Colors.white,
    }
});