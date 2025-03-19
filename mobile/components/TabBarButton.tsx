import { Colors } from "@/mobile/constants/Colors";
import { icon } from "@/mobile/constants/icons";
import { useCartStore } from "@/store/cardBadgeStore";

import { Pressable, StyleSheet, Text, View } from "react-native"


type Props = {
    onPress: () => void;
    onLongPress: () => void;
    isFocused: boolean;
    label: string;
    routeName: string
}

const TabBarButton = (props: Props) => {
    const { onPress, onLongPress, isFocused, label, routeName } = props;
    const { cartCount } = useCartStore();

    return (
        <Pressable
            onPress={onPress}
            onLongPress={onLongPress}
            style={styles.tabBarButton}
        >
            {(routeName == "cart" && cartCount > 0) && (
                < View style={styles.badgeWrapper}>
                    <Text style={styles.badge}>{cartCount}</Text>
                </View>
            )
            }
            {/* Cart Badge */}
            {icon[routeName]({ color: isFocused ? Colors.primary : Colors.black })}
            <Text style={{ color: isFocused ? Colors.secondary : Colors.darkGray }}>
                {label}
            </Text>
        </Pressable >
    );
}

export default TabBarButton

const styles = StyleSheet.create({
    tabBarButton: {
        flex: 1,
        alignItems: "center",
        justifyContent: "center",
        gap: 5,
    },
    badgeWrapper: {
        position: "absolute",
        backgroundColor: Colors.highlight,
        top: -5,
        right: 20,
        paddingVertical: 2,
        paddingHorizontal: 6,
        borderRadius: 10,
        zIndex: 10,
    },
    badge: {
        color: Colors.black,
        fontSize: 12,
    },
});