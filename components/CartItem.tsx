
import { View, Text, StyleSheet, Image, TouchableOpacity } from 'react-native'
import React from 'react'
import { CartItemType } from '@/types/type'
import { Colors } from '@/constants/Colors'
import { Ionicons } from '@expo/vector-icons'

type Props = {
    item: CartItemType
}


const CartItem = ({ item }: Props) => {
    return (
        <View style={styles.container}>
            <Image source={{ uri: item.image }} style={styles.itemImg} />
            <View style={styles.itemInfoWrapper}>
                <Text style={styles.itemTxt}>{item.title}</Text>
                <Text style={styles.itemTxt}>R${item.price}</Text>
                <View style={styles.itemControlWrapper}>
                    <TouchableOpacity>
                        <Ionicons name="trash-outline" size={20} color={"red"} />
                    </TouchableOpacity>
                    <View style={styles.quantityControlWrapper}>
                        <TouchableOpacity style={styles.quantityControl}>
                            <Ionicons name="remove-outline" size={20} color={Colors.black} />
                        </TouchableOpacity>
                        <Text style={styles.quantity}>1</Text>
                        <TouchableOpacity style={styles.quantityControl}>
                            <Ionicons name="add-outline" size={20} color={Colors.black} />
                        </TouchableOpacity>
                    </View>
                    <TouchableOpacity>
                        <Ionicons name="heart-outline" size={20} color={Colors.black} />
                    </TouchableOpacity>
                </View>
            </View>
        </View>
    )
}

export default CartItem

const styles = StyleSheet.create({
    container: {
        flexDirection: 'row',
        alignItems: 'center',
        padding: 10,
        marginBottom: 10,
        borderWidth: StyleSheet.hairlineWidth,
        borderColor: Colors.lightGray,
        borderRadius: 5,
        backgroundColor: Colors.white,
        width: "100%"
    },
    itemInfoWrapper: {
        flex: 1,
        alignSelf: "flex-start",
        gap: 10
    },
    itemImg: {
        width: 100,
        height: 100,
        borderRadius: 5,
        marginRight: 10
    },
    itemTxt: {
        fontSize: 16,
        fontWeight: '500',
        color: Colors.black,
        minWidth: 100,
        maxWidth: "80%",
    },
    itemControlWrapper: {
        flexDirection: 'row',
        justifyContent: 'space-between',
        alignItems: 'center',
        gap: 10,
    },
    quantityControlWrapper: {
        flexDirection: 'row',
        alignItems: 'center',
        gap: 15,
    },
    quantityControl:{
        padding: 5,
        borderWidth: 1,
        borderColor: Colors.lightGray,
        borderRadius: 5,
    },
    quantity: {
        fontSize: 16,
        fontWeight: '500',
        color: Colors.black
    }
});